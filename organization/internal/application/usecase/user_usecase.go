package usecase

import (
	"context"
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
