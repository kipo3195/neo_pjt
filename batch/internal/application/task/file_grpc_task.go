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
	GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error)
	ClearFileStorage(ctx context.Context, invalidFileIds []string, sendFileInfo map[string]string) error
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

func (r *fileGrpcTask) GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error) {
	return r.repository.GetInvalidFileInfo(ctx, yesterday)
}

func (r *fileGrpcTask) ClearFileStorage(ctx context.Context, invalidFileIds []string, sendFileInfo map[string]string) error {

	// sendFileInfo에 없는 file Id => 보내지지 않은 것 => 스토리지 삭제, error flag 변경
	clearFileId := make([]string, 0)
	// 실제로 보내졌지만 (message의 line_key 존재) send_flag가 업데이트 되지 않은 파일
	sendedFileId := make([]string, 0)
	for _, value := range invalidFileIds {
		if _, exists := sendFileInfo[value]; !exists {
			clearFileId = append(clearFileId, value)
		} else {
			sendedFileId = append(sendedFileId, value)
		}
	}
	log.Println("[ClearFileStorage] clear file :", clearFileId)
	log.Println("[ClearFileStorage] sendedFileId file :", sendedFileId)

	return r.repository.ClearFileStorage(ctx, clearFileId, sendedFileId)

}
