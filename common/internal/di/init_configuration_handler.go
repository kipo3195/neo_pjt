package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type ConfigurationHandler struct {
	Handler *handler.ConfigurationHandler
}

func InitConfigurationHandler(db *gorm.DB) *ConfigurationHandler {
	repository := repository.NewConfigurationRepository(db)
	usecase := usecase.NewConfigurationUsecase(repository, nil) // storage 필요함.
	handler := handler.NewConfigurationHandler(usecase)
	return &ConfigurationHandler{
		Handler: handler,
	}
}
