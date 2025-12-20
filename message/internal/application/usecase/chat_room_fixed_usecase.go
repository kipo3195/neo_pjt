package usecase

import "message/internal/domain/chatRoomFixed/repository"

type chatRoomFixedUsecase struct {
	repo repository.ChatRoomFixedRepository
}

type ChatRoomFixedUsecase interface {
}

func NewChatRoomFixedUsecase(repo repository.ChatRoomFixedRepository) ChatRoomFixedUsecase {
	return &chatRoomFixedUsecase{
		repo: repo,
	}
}
