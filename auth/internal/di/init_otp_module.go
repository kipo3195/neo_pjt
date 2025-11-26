package di

import (
	"auth/internal/application/usecase"
	"auth/internal/infrastructure/repository"
)

type OtpModule struct {
	Usecase usecase.OtpUsecase
}

func InitOtpModule(domain string) OtpModule {

	repo := repository.NewOtpApiRepository(domain)
	usecase := usecase.NewOtpUsecase(repo)

	return OtpModule{
		Usecase: usecase,
	}
}
