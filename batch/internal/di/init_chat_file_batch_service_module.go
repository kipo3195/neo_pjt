package di

import (
	"batch/internal/application/service"
	"batch/internal/application/task"
)

func InitChatFileBatchServiceModule(messageTask task.MessageGrpcTask, fileTask task.FileGrpcTask) service.ChatFileBatchService {
	return service.NewChatFileBatchService(messageTask, fileTask)
}
