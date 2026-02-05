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
}

func NewChatFileUsecase(repository repository.ChatFileRepository) ChatFileUsecase {
	return &chatFileUsecase{
		repository: repository,
	}
}

func (r *chatFileUsecase) UpDateFileStatus(ctx context.Context, transactionId string) error {
	return r.repository.UpdateFileStatus(ctx, transactionId)
}
