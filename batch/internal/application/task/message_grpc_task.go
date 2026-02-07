package task

import (
	"batch/internal/domain/messageGrpc/repository"
	"context"
	"log"
	"time"
)

type messageGrpcTask struct {
	repository repository.MessageGrpcRepository
}

type MessageGrpcTask interface {
	GetSendFileInfo(ctx context.Context, fileIds []string) error
}

func NewMessageGrpcTask(repository repository.MessageGrpcRepository) MessageGrpcTask {

	return &messageGrpcTask{
		repository: repository,
	}
}

func (r *messageGrpcTask) GetSendFileInfo(ctx context.Context, fileIds []string) error {

	checkDate := time.Now().AddDate(0, 0, -1)
	formattedDate := checkDate.Format("2006-01-02")

	log.Println("[sendFileInfo] check Date : ", formattedDate)

	r.repository.GetSendFileInfo(ctx, formattedDate, fileIds)
	return nil
}
