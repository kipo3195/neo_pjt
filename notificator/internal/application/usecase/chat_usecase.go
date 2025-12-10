package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/delivery/adapter"

	"notificator/internal/domain/chat/entity"
	"notificator/internal/infrastructure/storage"

	"github.com/gorilla/websocket"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	//repo repository.ChatRepository
	chatUserStorage storage.ChatUserStorage
}

type ChatUsecase interface {
	SubscribeChat(in input.ChatConnectInput, conn *websocket.Conn)
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput) output.ChatMessageOutput
}

func NewChatUsecase(chatUserStorage storage.ChatUserStorage) ChatUsecase {
	return &chatUsecase{
		chatUserStorage: chatUserStorage,
	}
}

func (u *chatUsecase) SubscribeChat(in input.ChatConnectInput, conn *websocket.Conn) {

	//entity := entity.MakeSubscribeChatEntity(in.UserHash)

	// 메모리에 사용자 정보 등록
	//u.chatUserStorage.PutChatConnect(entity.UserHash, conn)
}

// message broker가 아니더라도, rest api, rabbit mq를 통해 전달받은 데이터도 가공 처리 할 수 있다!
// 이게바로 클린 아키텍쳐!
// Input의 형태만 유지하면됨.
func (u *chatUsecase) RecvChatMessage(ctx context.Context, in input.ChatMessageInput) output.ChatMessageOutput {
	log.Println("[RecvChatMessage] recv data : ", in)

	chatLineEntity := entity.MakeChatLineEntity(in.ChatLineData.Cmd, in.ChatLineData.Contents, in.ChatLineData.LineKey, in.ChatLineData.TargetLineKey, in.ChatLineData.SendUserHash, in.ChatLineData.SendDate)
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoomData.RoomKey, in.ChatRoomData.RoomType, in.ChatRoomData.SecretFlag)
	en := entity.MakeRecvChatMessageEntity(in.EventType, in.ChatSession, chatRoomEntity, chatLineEntity)

	return adapter.MakeChatMessageOutput(en)

}
