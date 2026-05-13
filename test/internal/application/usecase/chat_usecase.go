package usecase

import (
	"context"
	"test/internal/application/usecase/input"
	"test/internal/infrastructure/config"
)

type chatUsecase struct {
	sfg *config.ServerConfig
}

type ChatUsecase interface {
	PutChatEvent(ctx context.Context, input input.PutChatInput)
}

func NewChatUsecase(sfg *config.ServerConfig) ChatUsecase {
	return &chatUsecase{
		sfg: sfg,
	}
}

func (r *chatUsecase) PutChatEvent(ctx context.Context, input input.PutChatInput) {

}
