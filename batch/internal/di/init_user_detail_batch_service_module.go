package di

import (
	"batch/internal/application/service"
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
)

func InitUserDetailBatchserviceModule(userDetail task.UserDetailTask, extendDbConnect task.ExtendDBConnectTask, serviceConfig *config.BatchConfig) service.UserDetailBatchService {

	return service.NewUserDetailBatchServiceModule(userDetail, extendDbConnect, serviceConfig)
}
