package usecase

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput) error
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn) ChatUsecase {
	return &chatUsecase{
		repository: repository,
		connector:  connector,
	}
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) error {

	chatLineEntity := entity.MakeChatLineEntity(in.ChatLine.Cmd, in.ChatLine.Contents, in.ChatLine.LineKey, in.ChatLine.SendUserHash, in.ChatLine.SendDate)
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoom.RoomKey, in.ChatRoom.RoomType)

	entity := entity.MakeSendChatEntity(in.EventType, in.ChatSession, chatLineEntity, chatRoomEntity)

	data, err := json.Marshal(entity) // 🔹 struct → []byte(JSON)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 채팅 발송
	err = u.connector.Publish("chat.message", data)
	if err != nil {
		log.Fatal("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError
	}

	return nil
}
