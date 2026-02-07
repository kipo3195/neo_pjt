package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/persistence/repository"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type FileGrpcModule struct {
	Task task.FileGrpcTask
}

func InitFileGrpcModule(db *gorm.DB, chatFileConfig config.ChatFileConfig, fileServiceGrpcClient *grpc.ClientConn) FileGrpcModule {

	repository := repository.NewFileGrpcRepository(db, fileServiceGrpcClient)
	task := task.NewFileGrpcTask(repository)

	return FileGrpcModule{
		Task: task,
	}
}
