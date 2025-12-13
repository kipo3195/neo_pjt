package di

import (
	"batch/internal/application/orchestrator"
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
)

func InitOrgInfoBatchServiceModule(orgInfo usecase.OrgInfoUsecase, extendDbConnect usecase.ExtendDBConnectUsecase, serviceConfig *config.OrgInfoBatchConfig) orchestrator.OrgInfoBatchService {

	return orchestrator.NewOrgInfoServiceModule(orgInfo, extendDbConnect, serviceConfig)
}
