package usecase

import (
	"context"
	"org/internal/delivery/dto/user"
	sharedEntity "org/internal/domain/shared/entity"
	"org/internal/domain/user/repository"
)

type userUsecase struct {
	repository repository.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, req user.GetMyInfoRequest) (user.GetMyInfoResponseBody, error)
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, req user.GetMyInfoRequest) (user.GetMyInfoResponseBody, error) {

	MyInfoEntity, err := r.repository.GetMyInfo(ctx, req.MyHash)
	if err != nil {
		return user.GetMyInfoResponseBody{}, err
	}
	return toMyInfoDto(MyInfoEntity), nil
}

func toMyInfoDto(entity sharedEntity.MyInfoEntity) user.GetMyInfoResponseBody {

	username := user.UsernameDto{
		Def: entity.Username.Ko, // 수정 필요
		Ko:  entity.Username.Ko,
		En:  entity.Username.En,
		Jp:  entity.Username.Jp,
		Zh:  entity.Username.Zh,
		Ru:  entity.Username.Ru,
		Vi:  entity.Username.Vi,
	}

	deptInfoDto := toDeptInfoDto(entity.DeptInfo)

	return user.GetMyInfoResponseBody{
		UserHash:     entity.UserHash,
		UserPhoneNum: entity.UserPhoneNum,
		Username:     username,
		OrgCodes:     []string{"neo"}, // TODO 수정 필요 메모리 기반, ACL
		ProfileUrl:   entity.ProfileUrl,
		ProfileMsg:   entity.ProfileMsg,
		DeptInfo:     deptInfoDto,
	}
}

func toDeptInfoDto(deptInfos []sharedEntity.DeptEntity) []user.DeptInfoDto {

	var deptInfoDto []user.DeptInfoDto

	for _, deptInfo := range deptInfos {
		deptInfoDto = append(deptInfoDto, user.DeptInfoDto{
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
