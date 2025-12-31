package sender

import (
	"context"
	"notificator/internal/domain/socketSender/entity"
)

type SocketSender interface {
	//SendChat(ctx context.Context, recv string, entity entity.SendChatEntity, conn *websocket.Conn) error
	SendChat(ctx context.Context, recv string, entity *entity.SendConnectionEntity, chatEntity entity.ChatEntity) error
	SendChatRoom(ctx context.Context, recv string, entity *entity.SendConnectionEntity, chatRoomEntity entity.ChatRoomEntity) error
}
