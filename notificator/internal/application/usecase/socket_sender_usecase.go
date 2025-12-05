package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"

	"github.com/gorilla/websocket"
)

type socketSenderUsecase struct {
	socketSender          sender.SocketSender
	sendConnectionStorage storage.SendConnectionStorage
}

type SocketSenderUsecase interface {
	SaveConnection(conn *websocket.Conn, userHash string)
	SendChat(ctx context.Context, input input.SendChatInput)
}

func NewSocketSenderUsecase(ss sender.SocketSender, sendConnectionStorage storage.SendConnectionStorage) SocketSenderUsecase {
	return &socketSenderUsecase{
		socketSender:          ss,
		sendConnectionStorage: sendConnectionStorage,
	}
}

func (r *socketSenderUsecase) SaveConnection(conn *websocket.Conn, userHash string) {

	// 쓰기용 채널 (Write Channel)
	Ch := make(chan interface{})
	entity := entity.MakeSendConnectionEntity(userHash, conn, Ch)
	r.sendConnectionStorage.PutConnection(userHash, entity)

	// 함수 종료시 세션 삭제
	defer r.sendConnectionStorage.RemoveConnection(userHash)

	for message := range entity.Chan {
		// 'ok' 검사가 필요 없습니다. 채널이 닫히면 for 루프는 자동으로 종료됩니다.

		// 이 위치에서 conn.WriteJSON 호출 (단일 고루틴에서만 접근하므로 안전)
		if err := entity.Conn.WriteJSON(message); err != nil {
			// 소켓이 닫혔을때 (클라이언트에서 닫으면, err -> websocket: close sent) 채널을 통해 데이터를 수신하지만, write할때 close send가 발생하므로
			//  고루틴 종료를 위해 break를 통해 종료 하여 루프를 빠져나가게 한다.
			log.Printf("[SaveConnection] userHash :%s, error : %s \n", userHash, err)
			break
		}
	}

	// 만약 쓰기 고루틴에서 에러가 발생하여 소켓을 닫을 일이 있다면 여기서 conn을 닫는것도 방법임.
}

func (r *socketSenderUsecase) SendChat(ctx context.Context, input input.SendChatInput) {

	sendChatRoomEntity := entity.MakeSendChatRoomEntity(input.ChatRoomData.RoomKey, input.ChatRoomData.RoomType, input.ChatRoomData.SecretFlag)
	sendChatLineEntity := entity.MakeSendChatLineEntity(input.ChatLineData.Cmd, input.ChatLineData.Contents, input.ChatLineData.LineKey, input.ChatLineData.SendUserHash, input.ChatLineData.SendDate)

	chatEntity := entity.MakeSendChatEntity(input.EventType, input.ChatSession, sendChatRoomEntity, sendChatLineEntity)

	// 이후 메모리에서 가져올 수 있도록 처리 필수
	RecvUserHash := make([]string, 0)
	RecvUserHash = append(RecvUserHash, "nauryhash", "kipo3195", "cyh8858hash")

	for _, recvUser := range RecvUserHash {

		connectionEntity := r.sendConnectionStorage.GetConnection(recvUser)

		if connectionEntity != nil {

			err := r.socketSender.SendChat(ctx, recvUser, connectionEntity, chatEntity)

			if err != nil {
				r.sendConnectionStorage.RemoveConnection(recvUser)
			}

		}
	}

	// for i := 0; i < len(RecvUserHash); i++ {

	// 	// 수신자의 웹소켓 connection 객체 조회

	// 	conn := r.chatUserStorage.GetChatConnect(RecvUserHash[i])

	// 	if conn == nil {
	// 		continue
	// 	}

	// 	err := r.socketSender.SendChat(ctx, RecvUserHash[i], entity, conn)

	// 	if err != nil {
	// 		r.chatUserStorage.RemoveChatConnect(RecvUserHash[i])
	// 	}
	// }
}
