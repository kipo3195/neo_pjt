package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type FileGrpcModule struct {
	Task task.FileGrpcTask
}

func InitFileGrpcModule(db *gorm.DB, chatFileConfig config.ChatFileConfig) FileGrpcModule {

	repository := repository.NewFileGrpcRepository(db)
	task := task.NewFileGrpcTask(repository)

	return FileGrpcModule{
		Task: task,
	}

}
