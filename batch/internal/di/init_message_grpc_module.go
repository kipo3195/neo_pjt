package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/persistence/repository"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type MessageGrpcModule struct {
	Task task.MessageGrpcTask
}

func InitMessageGrpcModule(db *gorm.DB, chatFileConfig config.ChatFileConfig, messageServiceGrpcClient *grpc.ClientConn) MessageGrpcModule {

	repository := repository.NewChatFileRepository(db, messageServiceGrpcClient)
	task := task.NewMessageGrpcTask(repository)

	return MessageGrpcModule{
		Task: task,
	}

}
