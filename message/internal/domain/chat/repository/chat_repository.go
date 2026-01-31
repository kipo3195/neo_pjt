package repository

import (
	"context"
	"message/internal/domain/chat/entity"
)

type ChatRepository interface {
	SaveChatLine(ctx context.Context, sendChatEntity entity.SendChatEntity) error
	ReadChatLine(ctx context.Context, readChatEntity entity.ReadChatEntity) error
	GetChatLineEvent(ctx context.Context, en entity.GetChatLineEventEntity) ([]entity.ChatLineEventEntity, error)
	GetChatFileEntity(ctx context.Context, transactionId string) ([]*entity.ChatFileEntity, error)
}
