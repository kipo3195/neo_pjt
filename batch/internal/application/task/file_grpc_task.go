package task

import (
	"batch/internal/domain/fileGrpc/repository"
	"context"
	"log"
	"time"
)

type fileGrpcTask struct {
	repository repository.FileGrpcRepository
}

type FileGrpcTask interface {
	UploadFileCheck(ctx context.Context) error
}

func NewFileGrpcTask(repository repository.FileGrpcRepository) FileGrpcTask {

	return &fileGrpcTask{
		repository: repository,
	}
}

func (r *fileGrpcTask) UploadFileCheck(ctx context.Context) error {

	checkDate := time.Now().AddDate(0, 0, -1)
	formattedDate := checkDate.Format("2006-01-02")

	log.Println("check Date : ", formattedDate)

	return r.repository.CheckUploadFile(ctx, formattedDate)
}
