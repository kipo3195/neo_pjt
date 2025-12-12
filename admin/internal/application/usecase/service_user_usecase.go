package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
	"admin/internal/delivery/adapter"
	"admin/internal/domain/serviceUser/entity"
	"admin/internal/domain/serviceUser/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
)

type serviceUserUsecase struct {
	repo repository.ServiceUserRepository
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

		if err != nil {
			log.Println("[RegistServiecUser] generateUserHash error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		temp := entity.ServiceUserEntity{
			UserId:   e,
			UserHash: hash,
		}

		serviceUsers = append(serviceUsers, temp)
	}

	result, err := r.repo.PutServiceUser(ctx, en.Org, serviceUsers)
	if err != nil {
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
