package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/consts"
	"notificator/internal/delivery/adapter"

	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/storage"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	repo                  repository.ChatRepository
	chatUserStorage       storage.ChatUserStorage
	sendConnectionStorage storage.SendConnectionStorage
}

type ChatUsecase interface {
	SubscribeChat(userHash string) error
	UnSubscribeChat(userHash string)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput) output.ChatMessageOutput
	RecvCreateChatRoomMessage(ctx context.Context, in input.CreateChatRoomMessageInput) error
}

func NewChatUsecase(chatUserStorage storage.ChatUserStorage, sendConnectionStorage storage.SendConnectionStorage, repo repository.ChatRepository) ChatUsecase {
	return &chatUsecase{
		chatUserStorage:       chatUserStorage,
		repo:                  repo,
		sendConnectionStorage: sendConnectionStorage,
	}
}

func (u *chatUsecase) SubscribeChat(userHash string) error {

	myChatRoom, err := u.repo.GetMyChatRoom(userHash)
	if err != nil {
		return err
	}

	u.chatUserStorage.InitMyRoom(myChatRoom.RoomKey, userHash)
	return nil
}

func (u *chatUsecase) UnSubscribeChat(userHash string) {
	u.chatUserStorage.CleanUpMyRoom(userHash)
}

// message broker가 아니더라도, rest api, rabbit mq를 통해 전달받은 데이터도 가공 처리 할 수 있다!
// 이게바로 클린 아키텍쳐!
// Input의 형태만 유지하면됨.
func (u *chatUsecase) RecvChatMessage(ctx context.Context, in input.ChatMessageInput) output.ChatMessageOutput {
	log.Println("[RecvChatMessage] recv data : ", in)

	chatLineEntity := entity.MakeChatLineEntity(in.ChatLineData.Cmd, in.ChatLineData.Contents, in.ChatLineData.LineKey, in.ChatLineData.TargetLineKey, in.ChatLineData.SendUserHash, in.ChatLineData.SendDate)
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoomData.RoomType, in.ChatRoomData.RoomKey, in.ChatRoomData.SecretFlag)
	en := entity.MakeRecvChatMessageEntity(in.EventType, in.ChatSession, chatRoomEntity, chatLineEntity)

	return adapter.MakeChatMessageOutput(en)

}

func (u *chatUsecase) RecvCreateChatRoomMessage(ctx context.Context, in input.CreateChatRoomMessageInput) error {

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

	if err != nil {
		return err
	}

	// sendConnectionStorage의 IsOnline이 true인 유저 = 실제 웹소켓 연결 유저만 별도로 조회
	onlineMember := make([]string, 0)
	for _, member := range chatRoomMemberEntity {

		isOnline := u.sendConnectionStorage.IsOnline(member.MemberHash)
		if isOnline {
			onlineMember = append(onlineMember, member.MemberHash)
		} else {
			log.Printf("[RecvCreateChatRoomMessage] userHash : %s is not connected. \n", member.MemberHash)
		}
	}
	u.chatUserStorage.PutChatRoomMember(createChatRoomEntity.RoomKey, onlineMember)

	return nil
}
