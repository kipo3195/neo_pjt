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

	// 별도의 고루틴으로 동작해야됨.
	// r.fileTask.UploadFileCheck(ctx)

	// ---------------------------
	// file service 에서 send_flag가 N인 데이터 가져오고

	yesterday := time.Now().AddDate(0, 0, -1)
	formattedDate := yesterday.Format("2006-01-02")

	invalidFileinfo, err := r.fileTask.GetInvalidFileInfo(ctx, formattedDate)
	if len(invalidFileinfo) == 0 {
		log.Println("[GetInvalidFileInfo] len = 0")
		log.Println("ChatFileBatchService end. time : ", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}

	sendFileInfo, err := r.messageTask.GetSendFileInfo(ctx, invalidFileinfo)
	if err != nil {
		return err
	}

	err = r.fileTask.ClearFileStorage(ctx, invalidFileinfo, sendFileInfo)
	if err != nil {
		return err
	}

	log.Println("ChatFileBatchService end. time : ", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
