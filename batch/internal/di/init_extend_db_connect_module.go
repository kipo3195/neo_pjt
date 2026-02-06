package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type ExtendDbConnectModule struct {
	Task task.ExtendDBConnectTask
}

func InitExtendDBConnectModule(db *gorm.DB) *ExtendDbConnectModule {

	repo := repository.NewExtendDBConnectRepository(db)
	task := task.NewExtendDBConnectTask(repo)

	return &ExtendDbConnectModule{
		Task: task,
	}
}
