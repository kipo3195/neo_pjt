package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type MessageGrpcModule struct {
	Task task.MessageGrpcTask
}

func InitMessageGrpcModule(db *gorm.DB, chatFileConfig config.ChatFileConfig) MessageGrpcModule {

	repository := repository.NewChatFileRepository(db)
	task := task.NewMessageGrpcTask(repository)

	return MessageGrpcModule{
		Task: task,
	}

}
