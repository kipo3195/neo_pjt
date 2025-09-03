package di

import (
	"core/internal/application/usecase"
	"core/internal/delivery/handler"
	"core/internal/infrastructure/config"
	"core/internal/infrastructure/repository"
	"core/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type AppValidationHandler struct {
	Handler *handler.AppValidationHandler
}

func InitAppValidationHandler(db *gorm.DB, sfg *config.ServerConfig, serverInfoStorage storage.ServerInfoStorage) *AppValidationHandler {

	repo := repository.NewAppValidationRepository(db)
	apiRepository := repository.NewAppValidationAPIRepository(sfg.Domain)
	usecase := usecase.NewAppValidationUsecase(repo, apiRepository, serverInfoStorage)
	handler := handler.NewAppValidationHandler(sfg, usecase)

	return &AppValidationHandler{
		Handler: handler,
	}

}
