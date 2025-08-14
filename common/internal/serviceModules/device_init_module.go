package serviceModules

import (
	configurationRepository "common/internal/domains/configuration/repositories/client"
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	"common/internal/domains/device/repositories/serverRepository"
	deviceUsecase "common/internal/domains/device/usecases/server"
	skinRepositories "common/internal/domains/skin/repositories/client"
	skinUsecase "common/internal/domains/skin/usecases/client"
	handlers "common/internal/serviceHandlers"
	"common/internal/services"
)

func InitDeviceInitModule(dep Dependencies) *handlers.DeviceInitHandler {

	deviceUsecase := deviceUsecase.NewDeviceUsecase(serverRepository.NewDeviceRepository(dep.DB))
	skinUsecase := skinUsecase.NewSkinUsecase(skinRepositories.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := configurationUsecase.NewConfigurationUsecase(configurationRepository.NewConfigurationRepository(dep.DB), dep.ConfigHashStorage)

	// 서비스 초기화
	svc := services.NewDeviceInitService(deviceUsecase, skinUsecase, configurationUsecase)

	// 핸들러 초기화
	return handlers.NewDeviceInitHandler(svc)
}
