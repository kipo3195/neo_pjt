package app

import (
	"org/internal/handler"
	"org/internal/infra/repository"
	"org/internal/usecase"

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
