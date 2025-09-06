package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type TokenHandler struct {
	Handler *handler.TokenHandler
}

func InitTokenHandler(db *gorm.DB, sfg *config.ServerConfig) *TokenHandler {

	repo := repository.NewTokenRepository(db)
	usecase := usecase.NewTokenUsecase(repo, sfg.GetJWTConfig())
	handler := handler.NewTokenHandler(usecase)

	return &TokenHandler{
		Handler: handler,
	}
}
