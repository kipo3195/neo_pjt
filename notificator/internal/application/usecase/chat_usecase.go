package usecase

import (
	"context"
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
