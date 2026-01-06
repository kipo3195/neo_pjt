package sender

import (
	"context"
	"notificator/internal/domain/socketSender/entity"
)

type ChatDataSender interface {
	SendChat(ctx context.Context, recv string, entity *entity.SendConnectionEntity, chatEntity entity.ChatEntity) error
}
