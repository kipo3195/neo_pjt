package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
	"common/internal/infrastructure/repository"
	"common/internal/services/dependencies"
)

type AppValidationHandler struct {
	Handler *handler.AppValidationHandler
}

func InitAppValidationService(dep dependencies.Dependency) *orchestrator.AppValidationService {

	appValidationUsecase := usecase.NewAppValidationUsecase(repository.NewAppValidationRepository(dep.DB), dep.ConfigHashStorage)
	skinUsecase := usecase.NewSkinUsecase(repository.NewSkinRepository(dep.DB), dep.SkinStorage)
	configurationUsecase := usecase.NewConfigurationUsecase(repository.NewConfigurationRepository(dep.DB), dep.ConfigHashStorage)

	// 핸들러 초기화
	return orchestrator.NewAppValidationService(appValidationUsecase, skinUsecase, configurationUsecase)
}
