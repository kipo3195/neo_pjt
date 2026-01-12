package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/consts"
	"notificator/internal/domain/chatRoom/entity"
	"notificator/internal/domain/chatRoom/repository"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/storage"
	"notificator/internal/util"

	"github.com/nats-io/nats.go"
)

type chatRoomUsecase struct {
	repo                  repository.ChatRoomRepository
	chatRoomStorage       storage.ChatRoomStorage
	connector             *nats.Conn
	messageSender         port.MessageSender
	sendConnectionStorage storage.SendConnectionStorage
}

type ChatRoomUsecase interface {
	SubscribeChat(userHash string) error
	UnSubscribeChat(userHash string)
	RecvCreateChatRoomMessage(ctx context.Context, in input.CreateChatRoomMessageInput) error
	RegistChatRoomMember(ctx context.Context, input input.ChatRoomEventInput)
	SendChatRoomEvent(ctx context.Context, input input.ChatRoomEventInput)
}

func NewChatRoomUsecase(repo repository.ChatRoomRepository, chatRoomStorage storage.ChatRoomStorage, sendConnectionStorage storage.SendConnectionStorage, connector *nats.Conn, messageSender port.MessageSender) ChatRoomUsecase {
	return &chatRoomUsecase{
		repo:                  repo,
		chatRoomStorage:       chatRoomStorage,
		connector:             connector,
		sendConnectionStorage: sendConnectionStorage,
		messageSender:         messageSender,
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
		consts.CHATROOM,
		consts.CREATE,
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
		log.Println(err)
		return err
	}

	// 여기서 다시 NATS로 발송 처리
	err = u.connector.Publish("chat.room.broadcast", data)
	if err != nil {
		log.Println("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError

		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}
	return nil
}

func (r *chatRoomUsecase) RegistChatRoomMember(ctx context.Context, input input.ChatRoomEventInput) {

	/* 참여자 entity 생성 */
	chatRoomMemberEntity := make([]string, 0)
	for _, m := range input.ChatRoomEventMemberInput {
		chatRoomMemberEntity = append(chatRoomMemberEntity, m.MemberHash)
	}

	/* 소켓 연결 참여자 선별 */
	socketConnectedMember := make([]string, 0)
	for _, m := range chatRoomMemberEntity {
		if r.sendConnectionStorage.GetConnection(m) != nil {
			socketConnectedMember = append(socketConnectedMember, m)
		}
	}

	// 참여자 regist
	r.chatRoomStorage.PutChatRoomMember(input.ChatRoomEventDataInput.RoomKey, socketConnectedMember)
}

func (r *chatRoomUsecase) SendChatRoomEvent(ctx context.Context, input input.ChatRoomEventInput) {

	chatRoomEntity := entity.MakeChatRoomEventEntity(input.Type, input.EventType, input.ChatRoomEventDataInput.CreateUserHash, input.ChatRoomEventDataInput.RegDate, input.ChatRoomEventDataInput.RoomKey, input.ChatRoomEventDataInput.RoomType, input.ChatRoomEventDataInput.Title, input.ChatRoomEventDataInput.SecretFlag, input.ChatRoomEventDataInput.Secret, input.ChatRoomEventDataInput.Description, input.ChatRoomEventDataInput.WorksCode, nil)

	/* 참여자 entity 생성 */
	chatRoomMemberEntity := r.chatRoomStorage.GetChatRoomMember(chatRoomEntity.ChatRoomData.RoomKey)

	chatRoomEventData := output.ChatRoomEventData{
		CreateUserHash: chatRoomEntity.ChatRoomData.CreateUserHash,
		RegDate:        chatRoomEntity.ChatRoomData.RegDate,
		RoomKey:        chatRoomEntity.ChatRoomData.RoomKey,
		RoomType:       chatRoomEntity.ChatRoomData.RoomType,
		Title:          chatRoomEntity.ChatRoomData.Title,
		SecretFlag:     chatRoomEntity.ChatRoomData.SecretFlag,
		Secret:         chatRoomEntity.ChatRoomData.Secret,
		Description:    chatRoomEntity.ChatRoomData.Description,
		WorksCode:      chatRoomEntity.ChatRoomData.WorksCode,
	}

	// dto 생성
	out := output.ChatRoomMessageOutput{
		Type:              consts.CHATROOM,
		EventType:         input.EventType,
		ChatRoomEventData: chatRoomEventData,
	}

	/* 웹소켓 연결된 사용자에게만 write */
	for _, recvUser := range chatRoomMemberEntity {
		r.messageSender.SendToClient(recvUser, out)
	}

}
