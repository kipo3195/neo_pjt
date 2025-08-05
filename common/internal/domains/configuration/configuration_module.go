package configuration

import "gorm.io/gorm"

type ConfigurationHandlers struct {
	ClientHandler *clientHandler.ConfigurationHandlers
}

func InitModule(db *gorm.DB) *ConfigurationHandlers {
	clientRepository := clientRepository.NewConfigurationRepository(db)
	clientUsecase := clientUsecase.NewConfigurationUsecase(clientRepository)
	clientHandler := clientHandler.NewConfigurationHandler(clientUsecase)

	return &ConfigurationHandlers{
		ClientHandler: clientHandler,
	}
}
