package sender

import (
	"context"
	"notificator/internal/consts"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/dto"
)

type webSocketSenderImpl struct {
}

func NewWebSocketSender() sender.SocketSender {
	return &webSocketSenderImpl{}
}

// func (r *webSocketSenderImpl) SendChat(ctx context.Context, recv string, entity entity.SendChatEntity, conn *websocket.Conn) error {
// 	// entity -> dto
// 	res := dto.MakeSendChatDto(entity.Type, entity.EventType, entity.ChatSession, entity.SendChatLineEntity, entity.SendChatRoomEntity)

// 	// 수신자 Hash 정보를 통해 websocket 객체를 storage에서 찾은 다음,
// 	// 해당 websocket에 write
// 	// Go 언어의 net/websocket 또는 가장 일반적으로 사용되는 gorilla/websocket 라이브러리의 websocket.Conn 객체는 동시 쓰기(Concurrent Write)에 안전하지 않습니다.
// 	// 즉, 여러 개의 고루틴이 하나의 conn 객체를 향해 동시에 WriteMessage, WriteJSON 등의 쓰기 메서드를 호출하면 데이터가 섞이거나(Corrupted),
// 	// 예상치 못한 오류가 발생할 수 있습니다 (경합 조건, Race Condition).
// 	if err := conn.WriteJSON(res); err != nil {
// 		log.Printf("websocket write error to %s: %v", recv, err)
// 		conn.Close()
// 		return err
// 	}

// 	return nil
// }

func (r *webSocketSenderImpl) SendChat(ctx context.Context, recv string, entity *entity.SendConnectionEntity, chatEntity entity.ChatEntity) error {

	res := dto.MakeSendChatResponse(chatEntity.Type, chatEntity.EventType, chatEntity.ChatSession, chatEntity.ChatLineEntity, chatEntity.ChatRoomEntity)

	select {
	case entity.Chan <- res:

	default:
		return consts.ErrSenderChannelError
	}

	return nil
}

func (r *webSocketSenderImpl) SendCreateChatRoom(ctx context.Context, recv string, entity *entity.SendConnectionEntity, en entity.CreateChatRoomEntity) error {

	res := dto.MakeCreateChatRoomResponse(
		en.CreateUserHash,
		en.RegDate,
		en.RoomKey,
		en.RoomType,
		en.Title,
		en.SecretFlag,
		en.Secret,
		en.Description,
		en.WorksCode,
	)

	select {
	case entity.Chan <- res:

	default:
		return consts.ErrSenderChannelError
	}

	return nil

}
