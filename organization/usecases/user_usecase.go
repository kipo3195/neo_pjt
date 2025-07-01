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

	MyInfoEntity, err := r.repo.GetMyInfo(ctx, req.MyHash)
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

func toDeptInfoDto(deptInfos []entities.DeptEntity) []userDto.DeptInfoDto {

	var deptInfoDto []userDto.DeptInfoDto

	for _, deptInfo := range deptInfos {
		deptInfoDto = append(deptInfoDto, userDto.DeptInfoDto{
			DeptOrg:  deptInfo.DeptOrg,
			DeptCode: deptInfo.DeptCode,
			DefLang:  deptInfo.DefLang,
			KoLang:   deptInfo.KoLang,
			EnLang:   deptInfo.EnLang,
			JpLang:   deptInfo.JpLang,
			ZhLang:   deptInfo.ZhLang,
			ViLang:   deptInfo.ViLang,
			RuLang:   deptInfo.RuLang,
			Header:   deptInfo.Header,
		})
	}

	return deptInfoDto
}
