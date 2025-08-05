package appToken

import (
	clientHandler "common/internal/domains/appToken/handlers/client"
	clientRepository "common/internal/domains/appToken/repositories/client"
	clientUsecase "common/internal/domains/appToken/usecases/client"

	"gorm.io/gorm"
)

type AppTokenHandlers struct {
	ClientHandler *clientHandler.AppTokenHandler
}

func InitModule(db *gorm.DB) *AppTokenHandlers {

	clientRepository := clientRepository.NewAppTokenRepository(db)
	clientUsecase := clientUsecase.NewAppTokenUsecase(clientRepository)
	clientHandler := clientHandler.NewAppTokenHandler(clientUsecase)

	return &AppTokenHandlers{
		ClientHandler: clientHandler,
	}
}
