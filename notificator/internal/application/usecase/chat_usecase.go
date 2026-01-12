package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/consts"
	"notificator/internal/domain/port"

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
	chatDebouncer   port.ChatCountDebouncer
}

type ChatUsecase interface {
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
	RecvChatCountMessage(ctx context.Context, in input.ChatCountMessageInput)
}

func NewChatUsecase(chatRoomStorage storage.ChatRoomStorage, repo repository.ChatRepository, messageSender port.MessageSender, chatDebouncer port.ChatCountDebouncer) ChatUsecase {
	return &chatUsecase{
		chatRoomStorage: chatRoomStorage,
		repo:            repo,
		messageSender:   messageSender,
		chatDebouncer:   chatDebouncer,
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

func (r *chatUsecase) RecvChatCountMessage(ctx context.Context, in input.ChatCountMessageInput) {

	chatCountEntity := entity.MakeChatCountEntity(in.RoomKey, in.RoomType, in.EventType, in.SendUserHash, in.Delta)
	log.Println("[RecvChatUnreadMessage] chatCountEntity: ", chatCountEntity)

	chatCountData := entity.ChatCountDataEntity{
		RoomKey:  chatCountEntity.RoomKey,
		RoomType: chatCountEntity.RoomType,
		Delta:    chatCountEntity.Delta,
	}

	chatCountMessageEntity := &entity.ChatCountMessageEntity{
		Type:          consts.CHATUNREAD,
		EventType:     in.EventType,
		ChatCountData: chatCountData,
	}

	if chatCountEntity.EventType == consts.READ {
		// 읽음처리 - 나에게 발송
		r.chatDebouncer.AddChatCount(chatCountEntity.SendUserHash, chatCountMessageEntity)
	} else if chatCountEntity.EventType == consts.UNREAD {
		// 신규 라인 발생 - 발신자를 제외하고 보냄.
		RecvUserHash := r.chatRoomStorage.GetChatRoomMember(in.RoomKey)
		for _, recvUser := range RecvUserHash {
			if recvUser != chatCountEntity.SendUserHash {
				r.chatDebouncer.AddChatCount(recvUser, chatCountMessageEntity)
			}
		}
	}

}
