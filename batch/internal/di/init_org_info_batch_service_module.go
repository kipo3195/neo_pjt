package di

import (
	"batch/internal/application/orchestrator"
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
)

func InitOrgInfoBatchServiceModule(orgInfo usecase.OrgInfoUsecase, extendDbConnect usecase.ExtendDBConnectUsecase, serviceConfig *config.BatchConfig) orchestrator.OrgInfoBatchService {

	return orchestrator.NewOrgInfoBatchServiceModule(orgInfo, extendDbConnect, serviceConfig)
}
