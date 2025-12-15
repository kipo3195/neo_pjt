package orchestrator

import (
	"common/internal/application/usecase"
)

type DeviceInitService struct {
	Device        usecase.WorksInfoUsecase
	Skin          usecase.SkinUsecase
	Configuration usecase.ConfigurationUsecase
	AppToken      usecase.AppTokenUsecase
	Org           usecase.OrgUsecase
}

func NewDeviceInitService(device usecase.WorksInfoUsecase, skin usecase.SkinUsecase, configuration usecase.ConfigurationUsecase, appToken usecase.AppTokenUsecase, org usecase.OrgUsecase) *DeviceInitService {
	return &DeviceInitService{device, skin, configuration, appToken, org}
}
