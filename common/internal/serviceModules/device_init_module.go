package serviceModules

import (
	appToken "common/internal/domains/appToken/usecases/server"
	configurationUsecase "common/internal/domains/configuration/usecases/server"
	"common/internal/domains/device/repositories/serverRepository"
	deviceUsecase "common/internal/domains/device/usecases/server"
	skinRepositories "common/internal/domains/skin/repositories/server"
	skinUsecase "common/internal/domains/skin/usecases/server"
	handlers "common/internal/serviceHandlers"
	"common/internal/services"
)

func InitDeviceInitModule(dep Dependencies) *handlers.DeviceInitHandler {

	deviceUsecase := deviceUsecase.NewDeviceUsecase(serverRepository.NewDeviceRepository(dep.DB))
	skinUsecase := skinUsecase.NewSkinUsecase(skinRepositories.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := configurationUsecase.NewConfigurationUsecase(dep.ConfigHashStorage)
	appTokenUsecase := appToken.NewAppTokenUsecase()

	// 서비스 초기화
	svc := services.NewDeviceInitService(deviceUsecase, skinUsecase, configurationUsecase, appTokenUsecase)

	// 핸들러 초기화
	return handlers.NewDeviceInitHandler(svc)
}
