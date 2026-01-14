package usecase

import (
	"context"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/delivery/util"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/workerPool"
	"time"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
	workerPool workerPool.ChatWorkerPool
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput) error
	addTaskChatLineJob(chatEntity entity.SendChatEntity, chatCountEventEntity entity.ChatCountEventEntity) error
	ReadChat(ctx context.Context, in input.ReadChatInput) error
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn, workerPool workerPool.ChatWorkerPool) ChatUsecase {
	// domain layer
	return &chatUsecase{
		repository: repository,
		connector:  connector,
		workerPool: workerPool, // Usecase는 ChatWorkerPool이라는 인터페이스에 의존하고, 이 인터페이스의 구현체가 chatWorkerPool 구조체라는 사실을 전혀 알지 못합니다.
	}
}
func (u *chatUsecase) ReadChat(ctx context.Context, in input.ReadChatInput) error {

	readChatEntity := entity.MakeReadChatEntity(in.RoomKey, in.RoomType, in.UserHash, in.ReadDate)

	err := u.repository.ReadChatLine(ctx, readChatEntity)

	if err != nil {
		log.Println(err)
		return err
	}

	chatCountEventEntity := entity.MakeChatCountEventEntity(readChatEntity.RoomKey, readChatEntity.RoomType, "read", readChatEntity.UserHash, 0)

	data, err := util.EntityMarshal(chatCountEventEntity)
	if err != nil {
		log.Println(err)
		return err
	}

	/* 미확인 건수 발송 Message Broker */
	err = u.connector.Publish("chat.count.broadcast", data)
	if err != nil {
		log.Println("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError
		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}

	return nil
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) error {

	// 채팅 라인 entity
	chatLineEntity := entity.MakeChatLineEntity(in.ChatLine.Cmd, in.ChatLine.Contents, in.ChatLine.LineKey, in.ChatLine.TargetLineKey, in.ChatLine.SendUserHash, in.ChatLine.SendDate)

	// 채팅 룸 entity
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoom.RoomKey, in.ChatRoom.RoomType, in.ChatRoom.SecretFlag)

	// 미확인 건수 전송용 entity
	chatCountEventEntity := entity.MakeChatCountEventEntity(in.ChatRoom.RoomKey, chatRoomEntity.RoomType, "unread", in.ChatLine.SendUserHash, 1)

	// 채팅 전송용 entity
	chatEntity := entity.MakeSendChatEntity(in.EventType, in.ChatSession, chatLineEntity, chatRoomEntity)

	log.Println("[SendChat] entity : ", chatEntity)

	data, err := util.EntityMarshal(chatEntity)
	if err != nil {
		log.Fatal(err)
		return err
	}

	/* 채팅 발송 Message Broker */
	err = u.connector.Publish("chat.broadcast", data)
	if err != nil {
		log.Println("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError

		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}

	/* DB 저장 Task */
	err = u.addTaskChatLineJob(chatEntity, chatCountEventEntity)
	return nil
}

func (u *chatUsecase) addTaskChatLineJob(entity entity.SendChatEntity, chatCountEventEntity entity.ChatCountEventEntity) error {

	// Context 전달
	// ※ http의 context 전달시 job 에서 repository 호출하는 DB 처리 프로세스에서 context canceled 발생.
	// 백그라운드 처리이므로 실제 http 요청에 대한 context가 아니기 때문
	// 그러므로 새로운 context 생성하여 전달하고 job의 호출이 끝난 시점에 cancel 호출로 수명 주기 관리 = AddTask 내부에서 defer cancle 필수!
	// 이후 다른 비동기 처리가 추가될때 context 를 밖에서 생성하고 각각의 task에 주입하므로써 하나의 context로 모든 task를 아우를수 있는지 점검해보기
	jobCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	job := &job.ChatLineJob{
		SendChatEntity:       entity,
		ChatCountEventEntity: chatCountEventEntity,
		Ctx:                  jobCtx,
		Cancel:               cancel,
		Connector:            u.connector,
	}

	u.workerPool.AddTask(job)
	return nil
}
