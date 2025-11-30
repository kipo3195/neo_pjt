package usecase

import (
	"message/internal/domain/chatRoom/repository"
	"message/internal/infrastructure/storage"
)

type chatRoomUsecase struct {
	repository repository.ChatRoomRepository
	storage    storage.ChatRoomStorage
}

type ChatRoomUsecase interface {
}

func NewChatRoomUsecase(repository repository.ChatRoomRepository, storage storage.ChatRoomStorage) ChatRoomUsecase {

	return &chatRoomUsecase{
		repository: repository,
		storage:    storage,
	}

}
