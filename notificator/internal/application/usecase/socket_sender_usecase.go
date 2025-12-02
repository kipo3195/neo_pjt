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

	for {
		select {
		case message, ok := <-entity.Chan:
			if !ok {
				// 채널이 닫혔을 경우 (서버 종료 또는 연결 정리)
				entity.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				// 소켓 닫혔을때
				log.Println("[SaveConnection] close connection userHash : ", userHash)
				return // 하지 않으면 for문 무한 루프
			}

			// 이 위치에서 conn.WriteJSON 호출 (단일 고루틴에서만 접근하므로 안전)
			if err := entity.Conn.WriteJSON(message); err != nil {
				// 쓰기 오류 발생 시 처리 (예: 연결 끊기)
				log.Println("[SaveConnection] write error : userHash", userHash)
				return
			}
			// case <- c.ctx.Done(): // Context 종료 신호 처리 (선택 사항)
			// 	return
		}
	}

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
