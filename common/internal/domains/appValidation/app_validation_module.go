package appvalidation

import (
	serverHandler "common/internal/domains/appValidation/handlers/server"
	serverRepository "common/internal/domains/appValidation/repositories/server"
	serverUsecase "common/internal/domains/appValidation/usecases/server"
	"common/internal/infra/storage"

	"gorm.io/gorm"
)

type AppValidationHandlers struct {
	ServerHandler *serverHandler.AppValidationHandler
}

func InitModule(db *gorm.DB, configHashStorage storage.ConfigHashStorage) *AppValidationHandlers {

	serverRepository := serverRepository.NewAppValidationRepository(db)
	serverUsecase := serverUsecase.NewAppValidationUsecase(serverRepository, configHashStorage)
	serverHandler := serverHandler.NewAppValidationHandler(serverUsecase)

	return &AppValidationHandlers{
		ServerHandler: serverHandler,
	}
}
