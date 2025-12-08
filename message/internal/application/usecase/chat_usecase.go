package usecase

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/pool"
	"message/internal/domain/chat/repository"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
	workerPool pool.ChatPool
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput) error
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn, workerPool pool.ChatPool) ChatUsecase {

	// domain layer
	return &chatUsecase{
		repository: repository,
		connector:  connector,
		workerPool: workerPool,
	}
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) error {

	// 채팅 라인 entity
	chatLineEntity := entity.MakeChatLineEntity(in.ChatLine.Cmd, in.ChatLine.Contents, in.ChatLine.LineKey, in.ChatLine.SendUserHash, in.ChatLine.SendDate)

	// 채팅 룸 entity
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoom.RoomKey, in.ChatRoom.RoomType, in.ChatRoom.SecretFlag)

	entity := entity.MakeSendChatEntity(in.EventType, in.ChatSession, chatLineEntity, chatRoomEntity)

	log.Println("[SendChat] send entity : ", entity)

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
