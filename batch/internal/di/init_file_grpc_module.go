package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/external/rpc"

	"google.golang.org/grpc"
)

type FileGrpcModule struct {
	Task task.FileGrpcTask
}

func InitFileGrpcModule(chatFileConfig config.ChatFileConfig, fileServiceGrpcClient *grpc.ClientConn) FileGrpcModule {

	// 필요시 db *gorm.DB 주입
	grpcRepository := rpc.NewFileGrpcRepository(fileServiceGrpcClient)
	task := task.NewFileGrpcTask(grpcRepository)

	return FileGrpcModule{
		Task: task,
	}
}
