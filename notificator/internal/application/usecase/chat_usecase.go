package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/delivery/dto/chat"

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
	SubscribeChat(chatMessage chat.ChatConnect, conn *websocket.Conn)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
}

func NewChatUsecase(chatUserStorage storage.ChatUserStorage) ChatUsecase {
	return &chatUsecase{
		chatUserStorage: chatUserStorage,
	}
}

func (u *chatUsecase) SubscribeChat(chatMessage chat.ChatConnect, conn *websocket.Conn) {

	// 메모리에 사용자 정보 등록
	u.chatUserStorage.PutChatConnect(chatMessage.UserHash, conn)
}

func (u *chatUsecase) RecvChatMessage(ctx context.Context, in input.ChatMessageInput) {
	log.Println("[RecvChatMessage] recv data : ", in)

	// 수신자 Hash 정보를 통해 websocket 객체를 storage에서 찾은 다음,
	// 해당 websocket에 write
	msg := map[string]string{
		"type":    in.Type,
		"from":    in.SendUserHash,
		"content": in.Contents,
		"lineKey": in.LineKey,
	}

	for i := 0; i < len(in.DestUserHash); i++ {

		conn := u.chatUserStorage.GetChatConnect(in.DestUserHash[i])

		if conn == nil {
			continue
		}

		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("websocket write error to %s: %v", in.DestUserHash[i], err)
			conn.Close()
			u.chatUserStorage.RemoveChatConnect(in.DestUserHash[i])
		}
	}

}
