package sender

import (
	"context"
	"log"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/dto"

	"github.com/gorilla/websocket"
)

type webSocketSenderImpl struct {
}

func NewWebSocketSender() sender.SocketSender {
	return &webSocketSenderImpl{}
}

func (r *webSocketSenderImpl) SendChat(ctx context.Context, recv string, entity entity.SendChatEntity, conn *websocket.Conn) error {
	// entity -> dto
	res := dto.MakeSendChatDto(entity.Type, entity.EventType, entity.ChatSession, entity.SendChatLineEntity, entity.SendChatRoomEntity)

	// 수신자 Hash 정보를 통해 websocket 객체를 storage에서 찾은 다음,
	// 해당 websocket에 write
	if err := conn.WriteJSON(res); err != nil {
		log.Printf("websocket write error to %s: %v", recv, err)
		conn.Close()
		return err
	}

	return nil
}
