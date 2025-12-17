package di

import (
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
	"admin/internal/infrastructure/repository"
)

type UserAuthRegisterModule struct {
	Usecase usecase.UserAuthRegisterUsecase
}

func InitUserAuthRegisterModule() UserAuthRegisterModule {

	repo := repository.NewUserAuthRegisterRepository()
	usecase := usecase.NewUserAuthRegisterUsecase(repo)
	_ = handler.NewUserAuthRegisterHandler(usecase)

	return UserAuthRegisterModule{
		Usecase: usecase,
	}

}
