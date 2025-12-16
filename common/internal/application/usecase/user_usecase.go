package usecase

import (
	"common/internal/application/usecase/input"
	"common/internal/application/usecase/output"
	"common/internal/consts"
	"common/internal/domain/user/entity"
	"common/internal/domain/user/repository"
	"common/internal/infrastructure/storage"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type userUsecase struct {
	repository    repository.UserRepository
	storage       storage.UserStorage
	apiRepository repository.UserAPIRepository
}

type UserUsecase interface {
	GenerateUserChallenge(ctx context.Context, input input.UserRegisterChallengeInput) (output.UserRegisterChallengeOutput, error)
	UserRegister(ctx context.Context, input input.UserRegisterInput) string
}

func NewUserUsecase(repository repository.UserRepository, storage storage.UserStorage, apiRepository repository.UserAPIRepository) UserUsecase {

	return &userUsecase{
		repository:    repository,
		storage:       storage,
		apiRepository: apiRepository,
	}
}

func (u userUsecase) GenerateUserChallenge(ctx context.Context, input input.UserRegisterChallengeInput) (output.UserRegisterChallengeOutput, error) {

	entity := entity.MakeUserRegisterChallengeEntity(input.Id, input.Salt)

	// DB체크 service_users에 user_id가 있는지 점검.

	err := u.repository.CheckUserRegist(ctx, entity.Id)

	if err != nil {
		return output.UserRegisterChallengeOutput{}, err
	}

	chanllenge := generateChallenge(entity.Id, entity.Salt)

	u.storage.PutUserChallenge(entity.Id, chanllenge)

	return output.UserRegisterChallengeOutput{
		Challenge: chanllenge,
	}, nil
}

func generateChallenge(userID, salt string) string {
	// 현재 시간 (밀리초 단위로 고유성 확보)
	now := time.Now().UnixNano()

	// 조합 문자열
	data := fmt.Sprintf("%s:%s:%d", userID, salt, now)

	// MD5 해싱
	hash := md5.Sum([]byte(data))

	// 32자리 문자열로 변환
	return hex.EncodeToString(hash[:])
}

func (u userUsecase) UserRegister(ctx context.Context, input input.UserRegisterInput) string {

	en := entity.MakeUserRegisterEntity(input.Id, input.Salt, input.Fv)

	hash := en.Fv
	challenge, err := u.storage.GetUserChallenge(en.Id)

	if err != nil {
		return consts.COMMON_F001
	}

	// CBC 복호화 호출
	plain, err := DecryptAES256CBC(hash, challenge)
	if err != nil {
		fmt.Println("Decrypt error:", err)
		return "code2"
	}

	var data map[string]interface{}
	if err := json.Unmarshal(plain, &data); err != nil {
		panic(err)
	}
	fmt.Println("map:", data)

	// ② struct 형태로
	var info entity.UserRegisterInfoEntity
	if err := json.Unmarshal(plain, &info); err != nil {
		fmt.Println("json unmarshal error")
		return "code2"
	}

	if input.Salt == info.Salt {
		u.storage.DeleteUserChallenge(en.Id)
		result, err := u.apiRepository.UserAuthRegistInAuth(ctx, en.Id, info, challenge)
		if result != "success" || err != nil {
			return "code2"
		}
		return "code1"
	} else {
		return "code2"
	}
}

// keyStr을 SHA-256으로 해시해서 32바이트 AES-256 키를 만듭니다.
func deriveKey(keyStr string) []byte {
	sum := sha256.Sum256([]byte(keyStr))
	return sum[:] // 32 bytes
}

func pkcs7Unpad(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("invalid padding size")
	}
	padLen := int(b[len(b)-1])
	if padLen == 0 || padLen > len(b) {
		return nil, errors.New("invalid padding")
	}
	// 패딩 값이 모두 padLen인지 확인
	for i := 0; i < padLen; i++ {
		if b[len(b)-1-i] != byte(padLen) {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return b[:len(b)-padLen], nil
}

func DecryptAES256CBC(base64Ciphertext, keyStr string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]         // 16
	ciphertext := data[aes.BlockSize:] // remainder

	key := deriveKey(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	plain, err = pkcs7Unpad(plain)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
