package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type OtpModule struct {
	Handler *handler.OtpHandler
}

func InitOtpModule(db *gorm.DB) *OtpModule {
	repository := repository.NewOtpApiRepository(db)
	usecase := usecase.NewOtpUsecase(repository)
	handler := handler.NewOtpHandler(usecase)
	return &OtpModule{
		Handler: handler,
	}
}
