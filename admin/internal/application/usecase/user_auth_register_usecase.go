package usecase

import (
	"admin/internal/domain/userAuthRegister/repository"
	"log"
)

type userAuthRegisterUsecase struct {
	repo    repository.UserAuthRegisterRepository
	apiRepo repository.UserAuthRegisterApiRepository
}

type UserAuthRegisterUsecase interface {
	UserAuthRegisterInAuth(ctx, interface{}) error
}

func NewUserAuthRegisterUsecase(repo repository.UserAuthRegisterRepository) UserAuthRegisterUsecase {

	return &userAuthRegisterUsecase{
		repo: repo,
	}
}

func (r *userAuthRegisterUsecase) UserAuthRegisterInAuth(ctx, interface{}) error {

	/* auth 서비스 호출 */
	err = r.apiRepo.UserAuthRegisterInAuth(ctx, serviceUsers)
	if err != nil {
		log.Println("[RegistServiecUser] auth server api call error")
		return err
	}

}
