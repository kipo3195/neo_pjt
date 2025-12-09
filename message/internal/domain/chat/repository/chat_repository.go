package repository

import (
	"context"
	"message/internal/domain/chat/entity"
)

type ChatRepository interface {
	SaveChatLine(ctx context.Context, sendChatEntity entity.SendChatEntity)
}
