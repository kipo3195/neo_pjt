package usecase

import "message/internal/domain/chatFile/repository"

type chatFileUsecase struct {
	repo repository.ChatFileRepository
}

type ChatFileUsecase interface {
}

func NewChatFileUsecase(repo repository.ChatFileRepository) ChatFileUsecase {
	return &chatFileUsecase{
		repo: repo,
	}
}
