package serviceModules

import (
	appValidationRepository "common/internal/domains/appValidation/repositories/server"
	appValidationUsecase "common/internal/domains/appValidation/usecases/server"
	configurationRepository "common/internal/domains/configuration/repositories/client"
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	skinRepositories "common/internal/domains/skin/repositories/client"
	skinUsecase "common/internal/domains/skin/usecases/client"
	"common/internal/services/dependencies"
	"common/internal/services/serviceDomains"
	"common/internal/services/serviceHandlers"
)

func InitAppValidationModule(dep dependencies.Dependency) *serviceHandlers.AppValidationHandler {

	appValidationUsecase := appValidationUsecase.NewAppValidationUsecase(appValidationRepository.NewAppValidationRepository(dep.DB), dep.ConfigHashStorage)
	skinUsecase := skinUsecase.NewSkinUsecase(skinRepositories.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := configurationUsecase.NewConfigurationUsecase(configurationRepository.NewConfigurationRepository(dep.DB), dep.ConfigHashStorage)

	// 서비스 초기화
	svc := serviceDomains.NewAppValidationService(appValidationUsecase, skinUsecase, configurationUsecase)

	// 핸들러 초기화
	return serviceHandlers.NewAppValidationHander(svc)
}
