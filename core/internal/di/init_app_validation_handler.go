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

	clientRepository := repository.NewAppValidationRepository(db)
	clientUsecase := usecase.NewAppValidationUsecase(clientRepository, serverInfoStorage)
	handler := handler.NewAppValidationHandler(sfg, clientUsecase)

	return &AppValidationHandler{
		Handler: handler,
	}

}
