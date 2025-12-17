package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
	"admin/internal/delivery/adapter"
	"admin/internal/domain/serviceUser/entity"
	"admin/internal/domain/serviceUser/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

type serviceUserUsecase struct {
	repo    repository.ServiceUserRepository
	apiRepo repository.ServiceUserApiRepository
}

type ServiceUserUsecase interface {
	RegistServiecUser(ctx context.Context, input input.RegistServiceUserInput) (output.RegistServiceUserOutput, error)
}

func NewServiceUserUsecase(repo repository.ServiceUserRepository) ServiceUserUsecase {
	return &serviceUserUsecase{
		repo: repo,
	}
}

func (r *serviceUserUsecase) RegistServiecUser(ctx context.Context, input input.RegistServiceUserInput) (output.RegistServiceUserOutput, error) {

	en := entity.MakeRegistServiceUserEntity(input.Org, input.UserId)

	serviceUsers := make([]entity.ServiceUserEntity, 0)

	for _, e := range en.UserId {

		hash, err := generateUserHash()

		// userHash 생성
		if err != nil {
			log.Println("[RegistServiecUser] generateUserHash error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		// 관리자를 통한 일괄 사용자 인증정보 등록 처리를 위한
		// salt 생성
		salt, err := generateSalt()
		if err != nil {
			log.Println("[RegistServiecUser] generateSalt error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		// 관리자를 통한 일괄 사용자 인증정보 등록 처리를 위한
		// userAuth 생성
		userAuth, err := generateUserAuth(input.UserAuth, salt)
		if err != nil {
			log.Println("[RegistServiecUser] generateUserAuth error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		temp := entity.ServiceUserEntity{
			UserId:   e,
			UserHash: hash,
			Salt:     salt,
			UserAuth: userAuth,
		}

		serviceUsers = append(serviceUsers, temp)
	}

	/* DB 저장 */
	result, err := r.repo.PutServiceUser(ctx, en.Org, serviceUsers)
	if err != nil {
		log.Println("[RegistServiecUser] db save error")
		return output.RegistServiceUserOutput{}, err
	}

	o := adapter.MakeRegistServiceUserOutput(result)

	return o, nil
}

func generateUserHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateUserAuth(hash string, salt string) (string, error) {
	const iterations = 150000
	const keyLen = 32

	dk := pbkdf2.Key([]byte(hash), []byte(salt), iterations, keyLen, sha256.New)

	// 저장 포맷: iterations:salt_base64:hash_base64
	saltB64 := base64.StdEncoding.EncodeToString([]byte(salt))
	hashB64 := base64.StdEncoding.EncodeToString(dk)

	return fmt.Sprintf("%d:%s:%s", iterations, saltB64, hashB64), nil
}
