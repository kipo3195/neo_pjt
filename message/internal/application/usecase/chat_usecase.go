package usecase

import (
	"context"
	"message/internal/application/usecase/input"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/broker"
)

type chatUsecase struct {
	repository repository.ChatRepository
	mb         broker.Broker
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput)
}

func NewChatUsecase(repository repository.ChatRepository, mb broker.Broker) ChatUsecase {
	return &chatUsecase{
		repository: repository,
		mb:         mb,
	}
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) {
	u.mb.PublishToChatUser("", in.Contents)
}
