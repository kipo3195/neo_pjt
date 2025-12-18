package di

import (
	"batch/internal/application/orchestrator"
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
)

func InitUserDetailBatchserviceModule(userDetail usecase.UserDetailUsecase, extendDbConnect usecase.ExtendDBConnectUsecase, serviceConfig *config.BatchConfig) orchestrator.UserDetailBatchService {

	return orchestrator.NewUserDetailBatchServiceModule(userDetail, extendDbConnect, serviceConfig)
}
