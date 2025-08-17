package server

import (
	deviceUsecase "common/internal/domains/device/usecases/server"
)

type DeviceHandler struct {
	usecase deviceUsecase.DeviceUsecase
}

func NewDeviceHandler(usecase deviceUsecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{
		usecase: usecase,
	}
}
