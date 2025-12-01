package usecase

import (
	"context"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"
)

type socketSenderUsecase struct {
	socketSender    sender.SocketSender
	chatUserStorage storage.ChatUserStorage
}

type SocketSenderUsecase interface {
	SendChat(ctx context.Context, input input.SendChatInput)
}

func NewSocketSenderUsecase(ss sender.SocketSender, chatUserStorage storage.ChatUserStorage) SocketSenderUsecase {
	return &socketSenderUsecase{
		socketSender:    ss,
		chatUserStorage: chatUserStorage,
	}
}

func (r *socketSenderUsecase) SendChat(ctx context.Context, input input.SendChatInput) {

	sendChatRoomEntity := entity.MakeSendChatRoomEntity(input.ChatRoomData.RoomKey, input.ChatRoomData.RoomType, input.ChatRoomData.SecretFlag)
	sendChatLineEntity := entity.MakeSendChatLineEntity(input.ChatLineData.Cmd, input.ChatLineData.Contents, input.ChatLineData.LineKey, input.ChatLineData.SendUserHash, input.ChatLineData.SendDate)

	entity := entity.MakeSendChatEntity(input.EventType, input.ChatSession, sendChatRoomEntity, sendChatLineEntity)

	// 이후 메모리에서 가져올 수 있도록 처리 필수
	RecvUserHash := make([]string, 0)
	RecvUserHash = append(RecvUserHash, "nauryhash", "kipo3195", "cyh8858hash")

	for i := 0; i < len(RecvUserHash); i++ {

		// 수신자의 웹소켓 connection 객체 조회
		conn := r.chatUserStorage.GetChatConnect(RecvUserHash[i])

		if conn == nil {
			continue
		}

		err := r.socketSender.SendChat(ctx, RecvUserHash[i], entity, conn)

		if err != nil {
			r.chatUserStorage.RemoveChatConnect(RecvUserHash[i])
		}
	}
}
