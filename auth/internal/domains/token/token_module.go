package token

import (
	"auth/internal/config"
	serverHandler "auth/internal/domains/token/handlers/server"
	serverRepository "auth/internal/domains/token/repositories/server"
	serverUsecase "auth/internal/domains/token/usecases/server"
	"auth/internal/utils"

	"gorm.io/gorm"
)

type TokenHandlers struct {
	//ClientHandler *clientHandler.TokenHnalder
	ServerHandler *serverHandler.TokenHandler
}

func InitModules(db *gorm.DB, jwtCfg *config.JWTConfig, authUtile *utils.AuthUtil) *TokenHandlers {

	serverRepository := serverRepository.NewTokenRepository(db)
	serverUsecase := serverUsecase.NewTokenUsecase(serverRepository, jwtCfg, authUtile)
	serverHandler := serverHandler.NewTokenHandler(serverUsecase)

	return &TokenHandlers{
		ServerHandler: serverHandler,
	}
}
