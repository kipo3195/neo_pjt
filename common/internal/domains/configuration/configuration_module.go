package configuration

import (
	clientHandler "common/internal/domains/configuration/handlers/client"
	clientRepository "common/internal/domains/configuration/repositories/client"
	clientUsecase "common/internal/domains/configuration/usecases/client"
	"common/internal/infra/storage"

	"gorm.io/gorm"
)

type ConfigurationHandlers struct {
	ClientHandler *clientHandler.ConfigurationHandler
}

func InitModule(db *gorm.DB, configStorage storage.ConfigHashStorage) *ConfigurationHandlers {
	clientRepository := clientRepository.NewConfigurationRepository(db)
	clientUsecase := clientUsecase.NewConfigurationUsecase(clientRepository, configStorage)
	clientHandler := clientHandler.NewConfigurationHandler(clientUsecase)

	return &ConfigurationHandlers{
		ClientHandler: clientHandler,
	}
}
