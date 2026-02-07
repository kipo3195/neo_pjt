package service

import (
	"batch/internal/application/task"
	"context"
	"log"
	"time"
)

type ChatFileBatchService struct {
	messageTask task.MessageGrpcTask
	fileTask    task.FileGrpcTask
}

func NewChatFileBatchService(messageTask task.MessageGrpcTask, fileTask task.FileGrpcTask) ChatFileBatchService {
	return ChatFileBatchService{
		messageTask: messageTask,
		fileTask:    fileTask,
	}
}

func (r *ChatFileBatchService) Run(ctx context.Context) error {
	log.Println("ChatFileBatchService start. time : ", time.Now().Format("2006-01-02 15:04:05"))

	r.fileTask.UploadFileCheck(ctx)

	// ---------------------------
	fileIds := make([]string, 0)
	fileIds = append(fileIds, "imgimg3")

	r.messageTask.GetSendFileInfo(ctx, fileIds)

	log.Println("ChatFileBatchService end. time : ", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
