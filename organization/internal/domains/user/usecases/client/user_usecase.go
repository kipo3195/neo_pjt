package client

import (
	"context"
	"org/entities"
	"org/internal/domains/user/dto/client/requestDTO"
	respositories "org/internal/domains/user/repositories/client"
)

type userUsecase struct {
	repository respositories.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, req requestDTO.GetMyInfoRequest) (requestDTO.GetMyInfoResponse, error)
}

func NewUserUsecase(repository respositories.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, req userDto.GetMyInfoRequest) (userDto.GetMyInfoResponse, error) {

	MyInfoEntity, err := r.repository.GetMyInfo(ctx, req.MyHash)
	if err != nil {
		return userDto.GetMyInfoResponse{}, err
	}
	return toMyInfoDto(MyInfoEntity), nil
}

func toMyInfoDto(entity entities.MyInfoEntity) userDto.GetMyInfoResponse {

	username := userDto.UsernameDto{
		Def: entity.Username.Ko, // 수정 필요
		Ko:  entity.Username.Ko,
		En:  entity.Username.En,
		Jp:  entity.Username.Jp,
		Zh:  entity.Username.Zh,
		Ru:  entity.Username.Ru,
		Vi:  entity.Username.Vi,
	}

	deptInfoDto := toDeptInfoDto(entity.DeptInfo)

	return userDto.GetMyInfoResponse{
		UserHash:     entity.UserHash,
		UserPhoneNum: entity.UserPhoneNum,
		Username:     username,
		OrgCodes:     []string{"neo"}, // TODO 수정 필요 메모리 기반, ACL
		ProfileUrl:   entity.ProfileUrl,
		ProfileMsg:   entity.ProfileMsg,
		DeptInfo:     deptInfoDto,
	}
}
