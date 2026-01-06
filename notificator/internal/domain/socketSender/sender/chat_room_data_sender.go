package sender

import (
	"context"
	"notificator/internal/domain/socketSender/entity"
)

type ChatRoomDataSender interface {
	SendCreateChatRoom(ctx context.Context, recv string, entity *entity.SendConnectionEntity, createChatRoomEntity entity.CreateChatRoomEntity) error
}
