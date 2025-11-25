package di

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase"
	"auth/internal/delivery/handler"
)

func InitDeviceAuthServiceModule(token usecase.TokenUsecase, device usecase.DeviceUsecase, otp usecase.OtpUsecase) *handler.DeviceAuthServiceHandler {
	service := orchestrator.NewDeviceAuthService(token, device, otp)
	return handler.NewDeviceAuthServiceHandler(service)
}
