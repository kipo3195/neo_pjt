package usecases

import (
	userDto "org/dto/client/user"
	"org/entities"
	"org/repositories"
)

type userUsecase struct {
	repo repositories.UserRepository
}

type UserUsecase interface {
}

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {

	return &userUsecase{
		repo: repo,
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
