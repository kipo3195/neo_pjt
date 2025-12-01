package sender

import (
	"context"
	"notificator/internal/domain/socketSender/entity"

	"github.com/gorilla/websocket"
)

type SocketSender interface {
	SendChat(ctx context.Context, recv string, entity entity.SendChatEntity, conn *websocket.Conn) error
}
