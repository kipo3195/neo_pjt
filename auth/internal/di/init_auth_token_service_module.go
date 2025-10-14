package di

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
)

func InitAuthTokenServiceModule(token usecase.TokenUsecase, userAuth usecase.UserAuthUsecase) *handler.UserAuthTokenServiceHandler {
	service := orchestrator.NewUserAuthTokenService(userAuth, token)
	return handler.NewUserAuthTokenServiceHandler(service)
}
