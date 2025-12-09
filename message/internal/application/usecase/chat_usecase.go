package usecase

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/workerPool"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
	workerPool workerPool.ChatWorkerPool
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput) error
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn, workerPool workerPool.ChatWorkerPool) ChatUsecase {
	// domain layer
	return &chatUsecase{
		repository: repository,
		connector:  connector,
		workerPool: workerPool, // Usecase는 ChatWorkerPool이라는 인터페이스에 의존하고, 이 인터페이스의 구현체가 chatWorkerPool 구조체라는 사실을 전혀 알지 못합니다.
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

	// 🎯 Job 생성 시 상위로부터 받은 Context를 Job에 담습니다.
	job := &job.ChatLineJob{
		SendChatEntity: entity,
		Ctx:            ctx, // Context 전달
	}

	u.workerPool.AddTask(job)

	return nil
}
