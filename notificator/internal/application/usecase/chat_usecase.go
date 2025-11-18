package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"

	"notificator/internal/domain/chat/entity"
	"notificator/internal/infrastructure/storage"

	"github.com/gorilla/websocket"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	//repo repository.ChatRepository
	chatUserStorage storage.ChatUserStorage
}

type ChatUsecase interface {
	SubscribeChat(in input.ChatConnectInput, conn *websocket.Conn)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
}

func NewChatUsecase(chatUserStorage storage.ChatUserStorage) ChatUsecase {
	return &chatUsecase{
		chatUserStorage: chatUserStorage,
	}
}

func (u *chatUsecase) SubscribeChat(in input.ChatConnectInput, conn *websocket.Conn) {

	entity := entity.MakeSubscribeChatEntity(in.UserHash)

	// 메모리에 사용자 정보 등록
	u.chatUserStorage.PutChatConnect(entity.UserHash, conn)
}

func (u *chatUsecase) RecvChatMessage(ctx context.Context, in input.ChatMessageInput) {
	log.Println("[RecvChatMessage] recv data : ", in)

	// 수신자 Hash 정보를 통해 websocket 객체를 storage에서 찾은 다음,
	// 해당 websocket에 write

	en := entity.MakeRecvChatMessageEntity(in.Type, in.LineKey, in.Contents, in.SendUserHash)

	for i := 0; i < len(in.RecvUserHash); i++ {

		// 수신자의 웹소켓 connection 객체 조회
		conn := u.chatUserStorage.GetChatConnect(in.RecvUserHash[i])

		if conn == nil {
			continue
		}

		if err := conn.WriteJSON(en); err != nil {
			log.Printf("websocket write error to %s: %v", in.RecvUserHash[i], err)
			conn.Close()
			u.chatUserStorage.RemoveChatConnect(in.RecvUserHash[i])
		}
	}

}
