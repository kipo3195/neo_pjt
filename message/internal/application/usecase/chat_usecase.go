package usecase

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput)
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn) ChatUsecase {
	return &chatUsecase{
		repository: repository,
		connector:  connector,
	}
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) {

	entity := entity.MakeSendChatEntity("", "", in.Contents, in.LineKey, nil)
	data, err := json.Marshal(entity) // 🔹 struct → []byte(JSON)
	if err != nil {
		log.Fatal(err)
	}

	// 채팅 발송
	err = u.connector.Publish("chat.message", data)
	if err != nil {
		log.Fatal("NATS publish failed:", err)
	}
}
