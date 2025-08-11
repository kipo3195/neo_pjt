package modules

import (
	appValidationRepository "common/internal/domains/appValidation/repositories/server"
	appValidationUsecase "common/internal/domains/appValidation/usecases/server"
	configurationRepository "common/internal/domains/configuration/repositories/client"
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	skinRepositories "common/internal/domains/skin/repositories/client"
	skinUsecase "common/internal/domains/skin/usecases/client"
	handlers "common/internal/handlers"
	"common/internal/infra/storage"
	"common/internal/services"

	"gorm.io/gorm"
)

type Dependencies struct {
	DB                *gorm.DB
	ConfigHashStorage storage.ConfigHashStorage
	SkinStorage       storage.SkinStorage
}

func InitAppInitModule(dep Dependencies) *handlers.AppInitHandler {

	appValidationUsecase := appValidationUsecase.NewAppValidationUsecase(appValidationRepository.NewAppValidationRepository(dep.DB), dep.ConfigHashStorage)
	skinUsecase := skinUsecase.NewSkinUsecase(skinRepositories.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := configurationUsecase.NewConfigurationUsecase(configurationRepository.NewConfigurationRepository(dep.DB), dep.ConfigHashStorage)

	// 서비스 초기화
	svc := services.NewAppInitService(appValidationUsecase, skinUsecase, configurationUsecase)

	// 핸들러 초기화
	return handlers.NewAppInitHander(svc)
}
