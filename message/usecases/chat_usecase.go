package usecases

import "message/repositories"

type chatUsecase struct {
	repo repositories.ChatRepository
}

type ChatUsecase interface {
}

func NewChatUsecase(repo repositories.ChatRepository) ChatUsecase {
	return &chatUsecase{repo: repo}
}
