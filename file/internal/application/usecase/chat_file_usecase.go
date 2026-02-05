package usecase

import "file/internal/domain/chatFile/repository"

type chatFileUsecase struct {
	repository repository.ChatFileRepository
}

type ChatFileUsecase interface {
}

func NewChatFileUsecase(repository repository.ChatFileRepository) ChatFileUsecase {
	return &chatFileUsecase{
		repository: repository,
	}
}
