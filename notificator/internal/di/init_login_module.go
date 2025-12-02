package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/delivery/handler"
	"notificator/internal/infrastructure/repository"

	"gorm.io/gorm"
)

type loginModule struct {
	Handler *handler.LoginHandler
	Usecase usecase.LoginUsecase
}

func InitLoginModule(db *gorm.DB) *loginModule {

	repository := repository.NewLoginRepository(db)
	usecase := usecase.NewLoginUsecase(repository)
	handler := handler.NewLoginHandler(usecase)

	return &loginModule{
		Handler: handler,
		Usecase: usecase,
	}
}
