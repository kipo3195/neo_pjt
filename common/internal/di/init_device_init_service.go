package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
)

func InitDeviceInitHandler(worksInfo usecase.WorksInfoUsecase, skin usecase.SkinUsecase, configuration usecase.ConfigurationUsecase, appToken usecase.AppTokenUsecase, org usecase.OrgUsecase) *handler.DeviceHandler {

	service := orchestrator.NewDeviceInitService(worksInfo, skin, configuration, appToken, org)
	return handler.NewDeviceHandler(service)
}
