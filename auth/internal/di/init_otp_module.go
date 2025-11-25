package di

import "auth/internal/application/usecase"

type OtpModule struct {
	Usecase usecase.OtpUsecase
}

func InitOtpModule() OtpModule {
	return OtpModule{}
}
