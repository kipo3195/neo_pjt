package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/consts"
	"auth/internal/domain/userAuth/entity"
	"auth/internal/domain/userAuth/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

type userAuthUsecase struct {
	repo    repository.UserAuthRepository
	storage storage.UserAuthStorage
}

type UserAuthUsecase interface {
	PutUserAuth(ctx context.Context, input input.UserAuthRegisterInput) string
	GenerateUserAuthChallenge(ctx context.Context, input input.UserAuthChallengeInput) (output.UserAuthChallengeOutput, error)
	GetUserAuth(ctx context.Context, input input.UserAuthInput) (output.UserAuthOutput, error)
}

func NewUserAuthUsecase(repo repository.UserAuthRepository, storage storage.UserAuthStorage) UserAuthUsecase {
	return userAuthUsecase{
		repo:    repo,
		storage: storage,
	}
}

func (u userAuthUsecase) PutUserAuth(ctx context.Context, input input.UserAuthRegisterInput) string {

	entity := entity.MakeUserAuthInfoEntity(input.Id, input.Salt, input.AuthHash, input.UserHash)

	err := u.repo.PutUserAuthInfo(ctx, entity)

	if err != nil {
		return "fail"
	}
	return "success"
}

func (u userAuthUsecase) GenerateUserAuthChallenge(ctx context.Context, input input.UserAuthChallengeInput) (output.UserAuthChallengeOutput, error) {

	entity := entity.MakeUserAuthChallengeEntity(input.Id)

	// DB 조회
	salt, err := u.repo.GetUserSalt(ctx, entity.Id)

	if err != nil {
		return output.UserAuthChallengeOutput{}, err
	}

	// challenge 생성
	challenge := generateChallenge(entity.Id, salt)

	// 메모리 저장 (TOBE redis)
	u.storage.PutUserAuthChallenge(entity.Id, challenge)

	return output.UserAuthChallengeOutput{
		Challenge: challenge,
		Salt:      salt,
	}, nil
}

func generateChallenge(userID string, salt string) string {
	// 현재 시간 (밀리초 단위로 고유성 확보)
	now := time.Now().UnixNano()

	// 조합 문자열
	data := fmt.Sprintf("%s:%s:%d", userID, salt, now)

	// SHA-256 해싱
	hash := sha256.Sum256([]byte(data))

	// 32자리 문자열로 변환 (원래 64자리)
	challenge := hex.EncodeToString(hash[:16])

	log.Printf("[generateChallenge] challenge : %s", challenge)

	return challenge
}

func (u userAuthUsecase) GetUserAuth(ctx context.Context, input input.UserAuthInput) (output.UserAuthOutput, error) {

	entity := entity.MakeUserAuthEntity(input.Id, input.Fv, input.Uuid)

	// id 기반으로 salt 찾기
	salt, err := u.repo.GetUserSalt(ctx, entity.Id)
	if err != nil {
		return output.UserAuthOutput{}, err
	}

	// id 기반으로 hash 찾기
	authHash, err := u.repo.GetUserAuthHash(ctx, entity.Id)
	if err != nil {
		return output.UserAuthOutput{}, err
	}

	// id 기반으로 challenge 찾기 (TOBE redis)
	challenge := u.storage.GetUserAuthChallenge(entity.Id)

	if challenge == "" {
		return output.UserAuthOutput{}, consts.ErrUserAuthChallengeExpired
	}

	hash := generateHash(authHash, salt)

	// 서버 hash + challenge
	serverFv := generateFV(hash, challenge)

	log.Println("GetUserAuth serverFn :", serverFv)
	log.Println("GetUserAuth entity.Fv :", entity.Fv)

	if serverFv == entity.Fv {
		// device를 기반으로 한 키가 있는지 여부 탐색 ... 일단은 없다 치고
		deviceChallenge := generateChallenge(entity.Uuid, challenge)
		return output.UserAuthOutput{
			DeviceChallenge: deviceChallenge,
		}, nil
	} else if serverFv != entity.Fv {
		// fv 불일치
		return output.UserAuthOutput{}, consts.ErrUserAuthFvMismatch
	}

	return output.UserAuthOutput{
		DeviceChallenge: "",
		AccessToken:     "",
		RefreshToken:    "",
	}, nil
}

func generateFV(hash string, challenge string) string {

	h := hmac.New(sha256.New, []byte(hash))
	h.Write([]byte(challenge))

	return hex.EncodeToString(h.Sum(nil))
}

func generateHash(hash string, salt string) string {
	const iterations = 150000
	const keyLen = 32

	dk := pbkdf2.Key([]byte(hash), []byte(salt), iterations, keyLen, sha256.New)

	// 저장 포맷: iterations:salt_base64:hash_base64
	saltB64 := base64.StdEncoding.EncodeToString([]byte(salt))
	hashB64 := base64.StdEncoding.EncodeToString(dk)

	return fmt.Sprintf("%d:%s:%s", iterations, saltB64, hashB64)
}
