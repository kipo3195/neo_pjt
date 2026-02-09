package usecase

import (
	"context"
	"file/internal/domain/chatFile/repository"
)

type chatFileUsecase struct {
	repository repository.ChatFileRepository
}

type ChatFileUsecase interface {
	UpDateFileStatus(ctx context.Context, transactionId string) error
	GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error)
	ClearFileStorage(ctx context.Context, clearFileIds []string, sendedFileIds []string) error
}

func NewChatFileUsecase(repository repository.ChatFileRepository) ChatFileUsecase {
	return &chatFileUsecase{
		repository: repository,
	}
}

func (r *chatFileUsecase) UpDateFileStatus(ctx context.Context, transactionId string) error {
	return r.repository.UpdateFileStatus(ctx, transactionId)
}

func (r *chatFileUsecase) GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error) {
	return r.repository.GetInvalidFileInfo(ctx, yesterday)
}

func (r *chatFileUsecase) ClearFileStorage(ctx context.Context, clearFileIds []string, sendedFileIds []string) error {

	// clearFileIds 스토리지 정리 -> TODO

	// sendedFileIds db 업데이트
	return r.repository.SendFlagUpdate(ctx, sendedFileIds)
}
