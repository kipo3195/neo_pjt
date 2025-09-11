package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
)

func InitDeviceInitHandler(worksInfo usecase.WorksInfoUsecase, skin usecase.SkinUsecase, configuration usecase.ConfigurationUsecase, appToken usecase.AppTokenUsecase) *handler.DeviceHandler {

	service := orchestrator.NewDeviceInitService(worksInfo, skin, configuration, appToken)

	return handler.NewDeviceHandler(service)
}
