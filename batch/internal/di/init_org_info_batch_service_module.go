package di

import (
	"batch/internal/application/service"
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
)

func InitOrgInfoBatchServiceModule(orgInfo task.OrgInfoTask, extendDbConnect task.ExtendDBConnectTask, serviceConfig *config.BatchConfig) service.OrgInfoBatchService {

	return service.NewOrgInfoBatchServiceModule(orgInfo, extendDbConnect, serviceConfig)
}
