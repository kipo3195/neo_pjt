package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/domain/userAuth/entity"
	"auth/internal/domain/userAuth/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type userAuthUsecase struct {
	repo    repository.UserAuthRepository
	storage storage.UserAuthStorage
}

type UserAuthUsecase interface {
	PutUserAuth(ctx context.Context, input input.UserAuthRegisterInput) string
	GenerateUserAuthChallenge(ctx context.Context, input input.UserAuthChallengeInput) (string, error)
	GetUserAuth(ctx context.Context, input input.UserAuthInput) output.UserAuthOutput
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

func (u userAuthUsecase) GenerateUserAuthChallenge(ctx context.Context, input input.UserAuthChallengeInput) (string, error) {

	// id 기반으로 salt를 조회해야할까?
	entity := entity.MakeUserAuthChallengeEntity(input.Id)

	salt, err := u.repo.GetUserSalt(ctx, entity.Id)

	if err != nil {
		return "", err
	}
	challenge := generateChallenge(entity.Id, salt)
	log.Printf("[GenerateUserAuthChallenge] challenge : %s", challenge)

	u.storage.PutUserAuthChallenge(entity.Id, challenge)
	return challenge, nil
}

func generateChallenge(userID string, salt string) string {
	// 현재 시간 (밀리초 단위로 고유성 확보)
	now := time.Now().UnixNano()

	// 조합 문자열
	data := fmt.Sprintf("%s:%s:%d", userID, salt, now)

	// MD5 해싱
	hash := md5.Sum([]byte(data))

	// 32자리 문자열로 변환
	return hex.EncodeToString(hash[:])
}

func (u userAuthUsecase) GetUserAuth(ctx context.Context, input input.UserAuthInput) output.UserAuthOutput {

	entity := entity.MakeUserAuthEntity(input.Id, input.Fv, input.Device)

	authHash, err := u.repo.GetUserAuthHash(ctx, entity.Id)
	if err != nil {
		return output.UserAuthOutput{
			DeviceChallenge: "",
			AccessToken:     "",
			RefreshToken:    "",
		}
	}

	challenge := u.storage.GetUserAuthChallenge(entity.Id)

	serverFv := generateHMAC(authHash, challenge)
	log.Println("GetUserAuth serverFn :", serverFv)
	if serverFv == entity.Fv {
		// device를 기반으로 한 키가 있는지 여부 탐색 ... 일단은 없다 치고
		deviceChallenge := generateChallenge(entity.Device, challenge)
		return output.UserAuthOutput{
			DeviceChallenge: deviceChallenge,
			AccessToken:     "",
			RefreshToken:    "",
		}
	}

	return output.UserAuthOutput{
		DeviceChallenge: "",
		AccessToken:     "",
		RefreshToken:    "",
	}
}

func generateHMAC(a string, b string) string {
	return ""
}
