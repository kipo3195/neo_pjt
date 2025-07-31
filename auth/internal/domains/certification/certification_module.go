package certification

import (
	"auth/internal/config"
	clientHandler "auth/internal/domains/certification/handlers/client"
	clientRepository "auth/internal/domains/certification/repositories/client"
	clientUsecase "auth/internal/domains/certification/usecases/client"
	"auth/internal/utils"

	"gorm.io/gorm"
)

type CertificationHandlers struct {
	ClientHandler *clientHandler.CertificationHandler
	//ServerHandler *serverHandler.CertificationHandler
}

func InitModules(db *gorm.DB, jwtCfg *config.JWTConfig, authUtile *utils.AuthUtil) *CertificationHandlers {
	clientRepository := clientRepository.NewCertificationRepository(db)
	clientUsecase := clientUsecase.NewCertificationUsecase(clientRepository, jwtCfg, authUtile)
	clientHandler := clientHandler.NewCertificationHandler(clientUsecase)

	return &CertificationHandlers{
		ClientHandler: clientHandler,
	}
}
