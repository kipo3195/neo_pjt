package usecase

import "message/internal/domain/chatRoomConfig/repository"

type chatRoomConfigUsecase struct {
	repo repository.ChatRoomConfigRepository
}

type ChatRoomConfigUsecase interface {
}

func NewChatRoomConfigUsecase(repo repository.ChatRoomConfigRepository) ChatRoomConfigUsecase {
	return &chatRoomConfigUsecase{
		repo: repo,
	}
}
