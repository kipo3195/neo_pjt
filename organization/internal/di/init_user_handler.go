package di

import (
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
	"org/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type UserHandler struct {
	Handler *handler.UserHandler
}

func InitUserHandler(db *gorm.DB) *UserHandler {

	repository := repository.NewUserRepository(db)
	usecase := usecase.NewUserUsecase(repository)
	handler := handler.NewUserHandler(usecase)

	return &UserHandler{
		Handler: handler,
	}
}
