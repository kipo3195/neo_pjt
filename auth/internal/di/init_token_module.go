package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type TokenModule struct {
	Handler *handler.TokenHandler
}

func InitTokenModule(db *gorm.DB, sfg *config.ServerConfig) *TokenModule {

	repo := repository.NewTokenRepository(db)
	usecase := usecase.NewTokenUsecase(repo, sfg.GetJWTConfig())
	handler := handler.NewTokenHandler(usecase)

	return &TokenModule{
		Handler: handler,
	}
}
