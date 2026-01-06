package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
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
	chatRoomStorage       storage.ChatRoomStorage
	sendConnectionStorage storage.SendConnectionStorage
}

type ChatUsecase interface {
	// SubscribeChat(userHash string) error
	// UnSubscribeChat(userHash string)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput) output.ChatMessageOutput
}

func NewChatUsecase(chatRoomStorage storage.ChatRoomStorage, sendConnectionStorage storage.SendConnectionStorage, repo repository.ChatRepository) ChatUsecase {
	return &chatUsecase{
		chatRoomStorage:       chatRoomStorage,
		repo:                  repo,
		sendConnectionStorage: sendConnectionStorage,
	}
}

// func (u *chatUsecase) SubscribeChat(userHash string) error {

// 	myChatRoom, err := u.repo.GetMyChatRoom(userHash)
// 	if err != nil {
// 		return err
// 	}

// 	u.chatRoomStorage.InitMyRoom(myChatRoom.RoomKey, userHash)
// 	return nil
// }

// func (u *chatUsecase) UnSubscribeChat(userHash string) {
// 	u.chatRoomStorage.CleanUpMyRoom(userHash)
// }

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
