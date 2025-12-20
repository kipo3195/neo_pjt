package usecase

import "message/internal/domain/chatRoomTitle/repository"

type chatRoomTitleUsecase struct {
	repo repository.ChatRoomTitleRepository
}

type ChatRoomTitleUsecase interface {
}

func NewChatRoomTitleUsecase(repo repository.ChatRoomTitleRepository) ChatRoomTitleUsecase {
	return &chatRoomTitleUsecase{
		repo: repo,
	}
}
