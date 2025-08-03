package skinImg

import (
	clientHandler "admin/internal/domains/skinImg/handlers/client"
	clientRepository "admin/internal/domains/skinImg/repositories/client"
	clientUsecase "admin/internal/domains/skinImg/usecases/client"

	"gorm.io/gorm"
)

type SkinImgHandlers struct {
	ClientHandler *clientHandler.SkinImgHandler
}

func InitModules(db *gorm.DB) *SkinImgHandlers {

	clientRepository := clientRepository.NewSkinImgRepository(db)
	clientUsecase := clientUsecase.NewSkinImgUsecase(clientRepository)
	clientHandler := clientHandler.NewSkinImgHandler(clientUsecase)

	return &SkinImgHandlers{
		ClientHandler: clientHandler,
	}
}
