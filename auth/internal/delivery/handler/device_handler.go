package handler

import "auth/internal/application/usecase"

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(uc usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: uc}
}
