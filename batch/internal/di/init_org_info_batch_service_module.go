package di

import (
	"batch/internal/application/orchestrator"
	"batch/internal/application/usecase"
)

func InitOrgInfoBatchServiceModule(orgInfo usecase.OrgInfoUsecase, extendDbConnect usecase.ExtendDBConnectUsecase) orchestrator.OrgInfoBatchService {

	service := orchestrator.NewOrgInfoServiceModule(orgInfo, extendDbConnect)
	return service
}
