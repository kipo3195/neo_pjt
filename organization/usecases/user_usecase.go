package usecases

import (
	"context"
	userDto "org/dto/client/user"
	"org/entities"
	"org/repositories"
)

type userUsecase struct {
	repo repositories.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, req userDto.GetMyInfoRequest) (userDto.GetMyInfoResponse, error)
}

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {

	return &userUsecase{
		repo: repo,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, req userDto.GetMyInfoRequest) (userDto.GetMyInfoResponse, error) {

	MyInfoEntity, err := r.repo.GetMyInfo(ctx, toGetMyInfoEntity(req))
	if err != nil {
		return userDto.GetMyInfoResponse{}, err
	}
	return toMyInfoDto(MyInfoEntity), nil
}

func toGetMyInfoEntity(req userDto.GetMyInfoRequest) entities.GetMyInfoEntity {
	return entities.GetMyInfoEntity{
		MyHash: req.MyHash,
	}
}

func toMyInfoDto(entity entities.MyInfoEntity) userDto.GetMyInfoResponse {

	username := userDto.UsernameDto{
		Kr: entity.Username.Kr,
		En: entity.Username.En,
		Jp: entity.Username.Jp,
		Cn: entity.Username.Cn,
	}

	return userDto.GetMyInfoResponse{
		UserHash:     entity.UserHash,
		UserPhoneNum: entity.UserPhoneNum,
		Username:     username,
	}
}
