package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type TokenModule struct {
	Handler *handler.TokenHandler
	Usecase usecase.TokenUsecase
}

func InitTokenModule(db *gorm.DB, sfg *config.ServerConfig, storage storage.AuthTokenStorage) *TokenModule {

	repo := repository.NewTokenRepository(db)
	usecase := usecase.NewTokenUsecase(repo, sfg.GetJWTConfig(), sfg.TokenConfig, storage)
	handler := handler.NewTokenHandler(usecase)

	return &TokenModule{
		Handler: handler,
		Usecase: usecase,
	}
}
