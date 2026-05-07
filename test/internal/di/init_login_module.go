package di

import (
	"test/internal/adapter/http/handler"
	"test/internal/application/usecase"
	"test/internal/infrastructure/config"
)

type LoginModule struct {
	Handler  *handler.LoginHandler
	Usercase usecase.LoginUsecase
}

func InitLoginModule(sfg *config.ServerConfig) *LoginModule {

	usecase := usecase.NewLoginUsecase(sfg)
	handler := handler.NewLoginHandler(usecase)

	return &LoginModule{
		Handler: handler,
	}
}
