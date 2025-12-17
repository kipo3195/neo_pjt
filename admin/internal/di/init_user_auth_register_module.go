package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"
)

type UserAuthRegisterModule struct {
	Usecase usecase.UserAuthRegisterUsecase
}

func InitUserAuthRegisterModule(domain string) UserAuthRegisterModule {

	repo := repository.NewUserAuthRegisterRepository()
	apiRepo := repository.NewServiceUserApiRepository(domain)
	usecase := usecase.NewUserAuthRegisterUsecase(repo, apiRepo)
	_ = handler.NewUserAuthRegisterHandler(usecase)

	return UserAuthRegisterModule{
		Usecase: usecase,
	}

}
