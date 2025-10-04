package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"org/internal/application/usecase/input"
	"org/internal/application/usecase/output"
	"org/internal/domain/user/entity"
	"org/internal/domain/user/repository"
)

type userUsecase struct {
	repository repository.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, input input.MyInfoInput) (output.MyInfoOutput, error)
	CreateServiceUser(ctx context.Context, input input.CreateServiceUserInput) error
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, input input.MyInfoInput) (output.MyInfoOutput, error) {

	entity := entity.MakeMyInfoHashEntity(input.MyHash)

	myInfo, err := r.repository.GetMyInfo(ctx, entity)

	if err != nil {
		return output.MyInfoOutput{}, err
	}
	output := output.MakeMyInfoOutput(myInfo)
	return output, nil
}

func (r *userUsecase) CreateServiceUser(ctx context.Context, input input.CreateServiceUserInput) error {

	var entities []entity.ServiceUserEntity
	for i := 0; i < input.UserCount; i++ {

		hash, err := generateUserHash()
		if err != nil {
			return fmt.Errorf("failed to generate user hash: %w", err)
		}

		userID := fmt.Sprintf("%s%04d", input.Keyword, i+1)

		entities = append(entities, entity.ServiceUserEntity{
			UserHash: hash,
			UserId:   userID,
		})
	}

	return r.repository.CreateServiceUser(ctx, entities)
}

func generateUserHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
