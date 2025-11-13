package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/delivery/dto/chat"

	"github.com/gorilla/websocket"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	//repo repository.ChatRepository
}

type ChatUsecase interface {
	SubscribeChat(chatMessage chat.ChatMessage, conn *websocket.Conn)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
}

func NewChatUsecase() ChatUsecase {
	return &chatUsecase{
		//repo: repo,
	}
}

func (u *chatUsecase) SubscribeChat(chatMessage chat.ChatMessage, conn *websocket.Conn) {

	// 메모리에 사용자 정보 등록

}

func (u *chatUsecase) RecvChatMessage(ctx context.Context, in input.ChatMessageInput) {
	log.Println("[RecvChatMessage] recv data : ", in)
}
