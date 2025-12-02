package usecase

import (
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/login/repository"
)

type loginUsecase struct {
	repository repository.LoginRepository
}

type LoginUsecase interface {
	LoginProcess(input input.LoginInput)
}

func NewLoginUsecase(repository repository.LoginRepository) LoginUsecase {
	return &loginUsecase{
		repository: repository,
	}
}

func (r *loginUsecase) LoginProcess(input input.LoginInput) {

	log.Printf("login api call uuid : %s, deviceType :%s \n", input.Uuid, input.DeviceType)
}
