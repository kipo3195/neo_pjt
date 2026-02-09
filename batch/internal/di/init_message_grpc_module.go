package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/external/rpc"

	"google.golang.org/grpc"
)

type MessageGrpcModule struct {
	Task task.MessageGrpcTask
}

func InitMessageGrpcModule(chatFileConfig config.ChatFileConfig, messageServiceGrpcClient *grpc.ClientConn) MessageGrpcModule {

	repository := rpc.NewChatFileRepository(messageServiceGrpcClient)
	task := task.NewMessageGrpcTask(repository)

	return MessageGrpcModule{
		Task: task,
	}

}
