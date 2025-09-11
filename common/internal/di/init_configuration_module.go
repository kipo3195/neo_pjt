package di

import (
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ConfigurationModule struct {
	Handler *handler.ConfigurationHandler
	Usecase usecase.ConfigurationUsecase
}

func InitConfigurationModule(db *gorm.DB, configHashStorage storage.ConfigHashStorage) *ConfigurationModule {
	repository := repository.NewConfigurationRepository(db)
	usecase := usecase.NewConfigurationUsecase(repository, configHashStorage) // storage 필요함.
	handler := handler.NewConfigurationHandler(usecase)
	return &ConfigurationModule{
		Handler: handler,
	}
}
