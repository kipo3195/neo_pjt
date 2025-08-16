package services

import (
	appToken "common/internal/domains/appToken/usecases/server"
	configuration "common/internal/domains/configuration/usecases/server"
	device "common/internal/domains/device/usecases/server"
	skin "common/internal/domains/skin/usecases/server"
)

type DeviceInitService struct {
	Device        device.DeviceUsecase
	Skin          skin.SkinUsecase
	Configuration configuration.ConfigurationUsecase
	AppToken      appToken.AppTokenUsecase
}

func NewDeviceInitService(device device.DeviceUsecase, skin skin.SkinUsecase, configuration configuration.ConfigurationUsecase, appToken appToken.AppTokenUsecase) *DeviceInitService {
	return &DeviceInitService{device, skin, configuration, appToken}
}
