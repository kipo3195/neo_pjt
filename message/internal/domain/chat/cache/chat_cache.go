package cache

import (
	"context"
	"message/internal/domain/chat/entity"
)

type ChatCache interface {
	GetFileEntity(ctx context.Context, transactionId string) ([]*entity.ChatFileEntity, error)
}
