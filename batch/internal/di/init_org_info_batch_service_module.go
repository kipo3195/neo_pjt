package di

import (
	orchestrator "batch/internal/application/service"
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
)

func InitOrgInfoBatchServiceModule(orgInfo task.OrgInfoTask, extendDbConnect task.ExtendDBConnectTask, serviceConfig *config.BatchConfig) orchestrator.OrgInfoBatchService {

	return orchestrator.NewOrgInfoBatchServiceModule(orgInfo, extendDbConnect, serviceConfig)
}
