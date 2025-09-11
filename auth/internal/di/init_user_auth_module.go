package di

import (
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type UserAuthModule struct {
	Handler *handler.UserAuthHandler
}

func InitUserAuthModule(db *gorm.DB, storage storage.UserAuthStorage) *UserAuthModule {

	repository := repository.NewUserAuthRepository(db)
	usecase := usecase.NewUserAuthUsecase(repository, storage)
	handler := handler.NewUserAuthHandler(usecase)

	return &UserAuthModule{
		Handler: handler,
	}
}
