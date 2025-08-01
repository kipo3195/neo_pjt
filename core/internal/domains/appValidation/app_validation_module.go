package appValidation

import (
	"gorm.io/gorm"

	"core/internal/config"
	clientHandler "core/internal/domains/appValidation/handlers/client"
	clientRepository "core/internal/domains/appValidation/repositories/client"
	clientUsecase "core/internal/domains/appValidation/usecases/client"
	"core/internal/infra/storage"
)

type AppValidationHandlers struct {
	ClientHandler *clientHandler.AppValidationHandler
}

func InitModules(db *gorm.DB, sfg *config.ServerConfig, serverInfoStorage storage.ServerInfoStorage) *AppValidationHandlers {

	clientRepository := clientRepository.NewAppValidationRepository(db)
	clientUsecase := clientUsecase.NewAppValidationUsecase(clientRepository, serverInfoStorage)
	clientHandler := clientHandler.NewAppValidationHandler(sfg, clientUsecase)

	return &AppValidationHandlers{
		ClientHandler: clientHandler,
	}

}
