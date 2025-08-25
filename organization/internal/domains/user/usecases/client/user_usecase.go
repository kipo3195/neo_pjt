package client

import (
	"context"
	"org/entities"
	"org/internal/domains/user/dto/client/requestDTO"
	"org/internal/domains/user/dto/client/responseDTO"
	respositories "org/internal/domains/user/repositories/client"
)

type userUsecase struct {
	repository respositories.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, req requestDTO.GetMyInfoRequest) (responseDTO.GetMyInfoResponseBody, error)
}

func NewUserUsecase(repository respositories.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, req requestDTO.GetMyInfoRequest) (responseDTO.GetMyInfoResponseBody, error) {

	MyInfoEntity, err := r.repository.GetMyInfo(ctx, req.MyHash)
	if err != nil {
		return responseDTO.GetMyInfoResponseBody{}, err
	}
	return toMyInfoDto(MyInfoEntity), nil
}

func toMyInfoDto(entity entities.MyInfoEntity) responseDTO.GetMyInfoResponseBody {

	username := responseDTO.UsernameDto{
		Def: entity.Username.Ko, // 수정 필요
		Ko:  entity.Username.Ko,
		En:  entity.Username.En,
		Jp:  entity.Username.Jp,
		Zh:  entity.Username.Zh,
		Ru:  entity.Username.Ru,
		Vi:  entity.Username.Vi,
	}

	deptInfoDto := toDeptInfoDto(entity.DeptInfo)

	return responseDTO.GetMyInfoResponseBody{
		UserHash:     entity.UserHash,
		UserPhoneNum: entity.UserPhoneNum,
		Username:     username,
		OrgCodes:     []string{"neo"}, // TODO 수정 필요 메모리 기반, ACL
		ProfileUrl:   entity.ProfileUrl,
		ProfileMsg:   entity.ProfileMsg,
		DeptInfo:     deptInfoDto,
	}
}

func toDeptInfoDto(deptInfos []entities.DeptEntity) []responseDTO.DeptInfoDto {

	var deptInfoDto []responseDTO.DeptInfoDto

	for _, deptInfo := range deptInfos {
		deptInfoDto = append(deptInfoDto, responseDTO.DeptInfoDto{
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
