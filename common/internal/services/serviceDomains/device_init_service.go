package servicesDomains

import (
	appToken "common/internal/domains/appToken/usecases/server"
	configuration "common/internal/domains/configuration/usecases/server"
	skin "common/internal/domains/skin/usecases/server"
	worksInfo "common/internal/domains/worksInfo/usecases/server"
)

type DeviceInitService struct {
	Device        worksInfo.WorksInfoUsecase
	Skin          skin.SkinUsecase
	Configuration configuration.ConfigurationUsecase
	AppToken      appToken.AppTokenUsecase
}

func NewDeviceInitService(device worksInfo.WorksInfoUsecase, skin skin.SkinUsecase, configuration configuration.ConfigurationUsecase, appToken appToken.AppTokenUsecase) *DeviceInitService {
	return &DeviceInitService{device, skin, configuration, appToken}
}
