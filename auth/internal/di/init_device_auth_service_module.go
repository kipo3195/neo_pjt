package di

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
)

func InitDeviceAuthServiceModule(token usecase.TokenUsecase, device usecase.DeviceUsecase) *handler.DeviceAuthServiceHandler {
	service := orchestrator.NewDeviceAuthService(token, device)
	return handler.NewDeviceAuthServiceHandler(service)
}
