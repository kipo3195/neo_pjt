package di

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type UserModule struct {
	Handler *handler.UserHandler
}

func InitUserModule(db *gorm.DB) *UserModule {

	repository := repository.NewUserRepository(db)
	usecase := usecase.NewUserUsecase(repository)
	handler := handler.NewUserHandler(usecase)

	return &UserModule{
		Handler: handler,
	}
}
