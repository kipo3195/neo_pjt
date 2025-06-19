package usecases

import (
	"context"
	userDto "org/dto/client/user"
	"org/repositories"
)

type userUsecase struct {
	repo repositories.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, req userDto.GetMyInfoRequest) (interface{}, error)
}

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {

	return &userUsecase{
		repo: repo,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, req userDto.GetMyInfoRequest) (interface{}, error) {
	return nil, nil
}
