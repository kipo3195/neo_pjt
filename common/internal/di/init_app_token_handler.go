package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type AppTokenHandler struct {
	Handler *handler.AppTokenHandler
}

func InitAppTokenHandler(db *gorm.DB) *AppTokenHandler {

	repository := repository.NewAppTokenRepository(db)
	usecase := usecase.NewAppTokenUsecase(repository)
	handler := handler.NewAppTokenHandler(usecase)

	return &AppTokenHandler{
		Handler: handler,
	}
}
