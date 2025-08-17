package serviceModules

import (
	appToken "common/internal/domains/appToken/usecases/server"
	configurationUsecase "common/internal/domains/configuration/usecases/server"
	skinRepositories "common/internal/domains/skin/repositories/server"
	skinUsecase "common/internal/domains/skin/usecases/server"
	"common/internal/domains/worksInfo/repositories/serverRepository"
	worksInfoUsecase "common/internal/domains/worksInfo/usecases/server"
	"common/internal/services/dependencies"
	"common/internal/services/serviceHandlers"
	"common/internal/services/servicesDomains"
)

func InitDeviceInitModule(dep dependencies.Dependency) *serviceHandlers.DeviceInitHandler {

	deviceUsecase := worksInfoUsecase.NewWorksInfoUsecase(serverRepository.NewWorksInfoRepository(dep.DB))
	skinUsecase := skinUsecase.NewSkinUsecase(skinRepositories.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := configurationUsecase.NewConfigurationUsecase(dep.ConfigHashStorage)
	appTokenUsecase := appToken.NewAppTokenUsecase()

	// 서비스 초기화
	svc := servicesDomains.NewDeviceInitService(deviceUsecase, skinUsecase, configurationUsecase, appTokenUsecase)

	// 핸들러 초기화
	return serviceHandlers.NewDeviceInitHandler(svc)
}
