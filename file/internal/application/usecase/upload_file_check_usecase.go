package usecase

import (
	"context"
	"file/internal/domain/uploadFileCheck/repository"
	"log"
)

type uploadFileCheckUsecase struct {
	repository repository.UploadFileCheckRepository
}

type UploadFileCheckUsecase interface {
	UploadFileCheck(ctx context.Context, checkDate string) error
}

func NewUploadFileCheckUsecase(repository repository.UploadFileCheckRepository) UploadFileCheckUsecase {
	return &uploadFileCheckUsecase{
		repository: repository,
	}
}

func (r *uploadFileCheckUsecase) UploadFileCheck(ctx context.Context, checkDate string) error {
	// 익일 스케쥴러로 체크시 전날의 upload_flag가 "N"인 것들을 조회해서 스토리지 확인 후 DB error_flag 변경하기

	invalidFile, err := r.repository.GetInvalidFile(ctx, checkDate)
	if err != nil {
		return err
	}

	log.Println("[UploadFileCheck] invalidFile : ", invalidFile)
	// 스토리지 접근, 있으면 삭제

	invalidFileId := make([]string, 0)
	for _, value := range invalidFile {
		invalidFileId = append(invalidFileId, value.FileId)
	}

	// error_flag 변경 -> message 서비스에서 send_chat 하더라도 발송 할 수 없는 상태로 만들어버림
	err = r.repository.UpdateInvalidFileState(ctx, invalidFileId)

	return nil
}
