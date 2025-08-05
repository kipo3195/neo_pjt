package skin

import (
	"gorm.io/gorm"

	serverHandler "common/internal/domains/skin/handlers/server"
	serverRepository "common/internal/domains/skin/repositories/server"
	serverUsecase "common/internal/domains/skin/usecases/server"

	clientHandler "common/internal/domains/skin/handlers/client"
	clientRepository "common/internal/domains/skin/repositories/client"
	clientUsecase "common/internal/domains/skin/usecases/client"
)

type SkinHandlers struct {
	ServerHandler *serverHandler.SkinHandler
	ClientHandler *clientHandler.SkinHandler
}

func InitModule(db *gorm.DB) *SkinHandlers {
	serverRepository := serverRepository.NewSkinRepository(db)
	serverUsecase := serverUsecase.NewSkinUsecase(serverRepository)
	serverHandler := serverHandler.NewSkinHandler(serverUsecase)

	clientRepository := clientRepository.NewSkinRepository(db)
	clientUsecase := clientUsecase.NewSkinUsecase(clientRepository)
	clientHandler := clientHandler.NewSkinHandler(clientUsecase)

	return &SkinHandlers{
		ServerHandler: serverHandler,
		ClientHandler: clientHandler,
	}
}
