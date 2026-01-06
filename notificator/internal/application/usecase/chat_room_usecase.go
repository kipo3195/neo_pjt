package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/consts"
	"notificator/internal/domain/chatRoom/entity"
	"notificator/internal/domain/chatRoom/repository"
	"notificator/internal/infrastructure/storage"
	"notificator/internal/util"

	"github.com/nats-io/nats.go"
)

type chatRoomUsecase struct {
	repo            repository.ChatRoomRepository
	chatRoomStorage storage.ChatRoomStorage
	connector       *nats.Conn
}

type ChatRoomUsecase interface {
	SubscribeChat(userHash string) error
	UnSubscribeChat(userHash string)
	RecvCreateChatRoomMessage(ctx context.Context, in input.CreateChatRoomMessageInput) error
}

func NewChatRoomUsecase(repo repository.ChatRoomRepository, chatRoomStorage storage.ChatRoomStorage, connector *nats.Conn) ChatRoomUsecase {
	return &chatRoomUsecase{
		repo:            repo,
		chatRoomStorage: chatRoomStorage,
		connector:       connector,
	}
}

func (u *chatRoomUsecase) SubscribeChat(userHash string) error {

	myChatRoom, err := u.repo.GetMyChatRoom(userHash)
	if err != nil {
		return err
	}

	u.chatRoomStorage.InitMyRoom(myChatRoom.RoomKey, userHash)
	return nil
}

func (u *chatRoomUsecase) UnSubscribeChat(userHash string) {
	u.chatRoomStorage.CleanUpMyRoom(userHash)
}

func (u *chatRoomUsecase) RecvCreateChatRoomMessage(ctx context.Context, in input.CreateChatRoomMessageInput) error {

	chatRoomMemberEntity := make([]entity.ChatRoomMemberEntity, 0)

	for _, v := range in.CreateChatRoomMemberInput {

		temp := entity.ChatRoomMemberEntity{
			MemberHash:      v.MemberHash,
			MemberWorksCode: v.MemberWorksCode,
		}

		chatRoomMemberEntity = append(chatRoomMemberEntity, temp)
	}

	if len(chatRoomMemberEntity) == 0 {
		return consts.ErrChatRoomMemberInvalid
	}

	createChatRoomEntity := entity.MakeCreateChatRoomEntity(in.CreateChatRoomInput.RoomKey, in.CreateChatRoomInput.RoomType, chatRoomMemberEntity)
	err := u.repo.PutChatRoomMember(ctx, createChatRoomEntity)

	chatRoomEventEntity := entity.MakeChatRoomEventEntity(
		"chatRoom",
		"C",
		in.CreateChatRoomInput.CreateUserHash,
		in.CreateChatRoomInput.RegDate,
		in.CreateChatRoomInput.RoomKey,
		in.CreateChatRoomInput.RoomType,
		in.CreateChatRoomInput.Title,
		in.CreateChatRoomInput.SecretFlag,
		in.CreateChatRoomInput.Secret,
		in.CreateChatRoomInput.Description,
		in.CreateChatRoomInput.WorksCode,
		chatRoomMemberEntity,
	)

	if err != nil {
		return err
	}

	data, err := util.EntityMarshal(chatRoomEventEntity)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 여기서 다시 NATS로 발송 처리
	err = u.connector.Publish("chat.room.broadcast", data)
	if err != nil {
		log.Fatal("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError

		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}
	return nil
}
