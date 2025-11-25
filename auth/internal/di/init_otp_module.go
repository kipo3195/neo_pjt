package di

import (
	"auth/internal/application/usecase"
	"auth/internal/infrastructure/repository"
)

type OtpModule struct {
	Usecase usecase.OtpUsecase
}

func InitOtpModule() OtpModule {

	repo := repository.NewOtpApiRepository()
	usecase := usecase.NewOtpUsecase(repo)

	return OtpModule{
		Usecase: usecase,
	}
}
