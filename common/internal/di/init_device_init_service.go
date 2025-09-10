package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/infrastructure/repository"
	"common/internal/services/dependencies"
)

func InitDeviceInitService(dep dependencies.Dependency) *orchestrator.DeviceInitService {

	worksInfoUsecase := usecase.NewWorksInfoUsecase(repository.NewWorksInfoRepository(dep.DB))
	skinUsecase := usecase.NewSkinUsecase(repository.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := usecase.NewConfigurationUsecase(repository.NewConfigurationRepository(dep.DB), dep.ConfigHashStorage)
	appTokenUsecase := usecase.NewAppTokenUsecase(repository.NewAppTokenRepository(dep.DB))

	// 서비스 초기화 초기화
	return orchestrator.NewDeviceInitService(worksInfoUsecase, skinUsecase, configurationUsecase, appTokenUsecase)
}
