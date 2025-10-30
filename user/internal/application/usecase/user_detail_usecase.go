package usecase

import (
	"context"
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/delivery/adapter"
	"user/internal/domain/userDetail/entity"
	"user/internal/domain/userDetail/repository"
)

type userDetailUsecase struct {
	repository repository.UserDetailRepository
}

type UserDetailUsecase interface {
	GetUserDetailInfo(ctx context.Context, input input.GetUserDetailInfoInput) (output.GetUserDetailInfoOutput, error)
}

func NewUserDatailUsecase(repository repository.UserDetailRepository) UserDetailUsecase {
	return &userDetailUsecase{
		repository: repository,
	}
}

func (u *userDetailUsecase) GetUserDetailInfo(ctx context.Context, input input.GetUserDetailInfoInput) (output.GetUserDetailInfoOutput, error) {

	entity := entity.MakeGetUserDetailInfoEntity(input.UserHashs)

	userInfos, err := u.repository.GetUserInfoDetailInfo(ctx, entity)

	if err != nil {
		return output.GetUserDetailInfoOutput{}, err
	}

	output := adapter.MakeGetUserDetailInfoOutput(userInfos)

	return output, nil
}
