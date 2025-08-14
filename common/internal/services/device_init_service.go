package services

import (
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	deviceUsecase "common/internal/domains/device/usecases/server"
	skinUsecase "common/internal/domains/skin/usecases/client"
)

type DeviceInitService struct {
	Device        deviceUsecase.DeviceUsecase
	Skin          skinUsecase.SkinUsecase
	Configuration configurationUsecase.ConfigurationUsecase
}

func NewDeviceInitService(v deviceUsecase.DeviceUsecase, s skinUsecase.SkinUsecase, c configurationUsecase.ConfigurationUsecase) *DeviceInitService {
	return &DeviceInitService{v, s, c}
}
