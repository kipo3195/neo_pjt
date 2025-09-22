package di

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
)

func InitUserAuthServiceModule(userAuth usecase.UserAuthUsecase, device usecase.DeviceUsecase) *handler.UserAuthServiceHandler {
	service := orchestrator.NewUserAuthService(userAuth, device)
	return handler.NewUserAuthServiceHandler(service)
}
