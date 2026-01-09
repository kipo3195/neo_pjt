package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/consts"
	"notificator/internal/core/port"

	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/storage"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	repo            repository.ChatRepository
	chatRoomStorage storage.ChatRoomStorage
	messageSender   port.MessageSender
}

type ChatUsecase interface {
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
	RecvChatUnreadMessage(ctx context.Context, in input.ChatUnreadMessageInput)
}

func NewChatUsecase(chatRoomStorage storage.ChatRoomStorage, repo repository.ChatRepository, messageSender port.MessageSender) ChatUsecase {
	return &chatUsecase{
		chatRoomStorage: chatRoomStorage,
		repo:            repo,
		messageSender:   messageSender,
	}
}

func (r *chatUsecase) RecvChatMessage(ctx context.Context, input input.ChatMessageInput) {

	chatRoomEntity := entity.MakeChatRoomEntity(input.ChatRoomData.RoomKey, input.ChatRoomData.RoomType, input.ChatRoomData.SecretFlag)
	chatLineEntity := entity.MakeChatLineEntity(input.ChatLineData.Cmd, input.ChatLineData.Contents, input.ChatLineData.LineKey, input.ChatLineData.TargetLineKey, input.ChatLineData.SendUserHash, input.ChatLineData.SendDate)

	RecvUserHash := r.chatRoomStorage.GetChatRoomMember(input.ChatRoomData.RoomKey)

	chatRoomOutput := output.ChatRoomDataOutput{
		RoomKey:    chatRoomEntity.RoomKey,
		RoomType:   chatRoomEntity.RoomType,
		SecretFlag: chatRoomEntity.SecretFlag,
	}

	chatLintOutput := output.ChatLineDataOutput{
		Cmd:           chatLineEntity.Cmd,
		Contents:      chatLineEntity.Contents,
		LineKey:       chatLineEntity.LineKey,
		TargetLineKey: chatLineEntity.TargetLineKey,
		SendUserHash:  chatLineEntity.SendUserHash,
		SendDate:      chatLineEntity.SendDate,
	}

	out := output.ChatMessageOutput{
		Type:         consts.CHAT,
		EventType:    input.EventType,
		ChatSession:  input.ChatSession,
		ChatRoomData: chatRoomOutput,
		ChatLineData: chatLintOutput,
	}

	for _, recvUser := range RecvUserHash {
		r.messageSender.SendToClient(recvUser, out)
	}
}

func (r *chatUsecase) RecvChatUnreadMessage(ctx context.Context, in input.ChatUnreadMessageInput) {

	chatUnreadEntity := entity.MakeChatUnreadEntity(in.RoomKey, in.RoomType, in.UnreadType, in.SendUserHash, in.Delta)
	log.Println("[RecvChatUnreadMessage] chatUnreadEntity: ", chatUnreadEntity)

	RecvUserHash := r.chatRoomStorage.GetChatRoomMember(in.RoomKey)

	chatUnreadOutput := output.ChatUnreadDataOutput{
		RoomKey:  chatUnreadEntity.RoomKey,
		RoomType: chatUnreadEntity.RoomType,
		Delta:    chatUnreadEntity.Delta,
	}

	out := output.ChatUnreadMessageOutput{
		Type:           consts.CHATUNREAD,
		EventType:      in.UnreadType,
		ChatUnreadData: chatUnreadOutput,
	}

	// TODO 사용자별 buffer 처리 로직 추가 필요.
	// 2~3개 이상 쌓였거나 대기시간이 0.5초 이상이거나 읽음처리가 와서 0으로 변경해야 하는 경우에 발송

	// 발신자를 제외하고 보냄.
	for _, recvUser := range RecvUserHash {
		//if recvUser != chatUnreadEntity.SendUserHash {
		r.messageSender.SendToClient(recvUser, out)
		//}
	}

}
