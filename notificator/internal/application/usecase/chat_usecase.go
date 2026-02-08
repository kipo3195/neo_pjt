package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/application/usecase/output"
	"notificator/internal/consts"
	"notificator/internal/domain/port"
	"notificator/pkg/dto"

	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/persistence/storage"
	"notificator/internal/infrastructure/workerPool"
)

const (
	CMD = "cmd" // 이벤트 type의 세부 타입
)

type chatUsecase struct {
	repo                   repository.ChatRepository
	chatRoomStorage        storage.ChatRoomStorage
	messageSender          port.MessageSender
	chatCountWorkerPool    workerPool.ChatCountWorkerPool
	chatReadDateWorkerPool workerPool.ChatReadDateWorkerPool
}

type ChatUsecase interface {
	RecvChatMessage(ctx context.Context, in input.ChatMessageInput)
	RecvChatCountMessage(ctx context.Context, in input.ChatCountMessageInput)
	RecvChatReadMessage(ctx context.Context, in input.ChatReadMessageInput)
}

func NewChatUsecase(chatRoomStorage storage.ChatRoomStorage, repo repository.ChatRepository, messageSender port.MessageSender, chatCountWorkerPool workerPool.ChatCountWorkerPool, chatReadDateWorkerPool workerPool.ChatReadDateWorkerPool) ChatUsecase {
	return &chatUsecase{
		chatRoomStorage:        chatRoomStorage,
		repo:                   repo,
		messageSender:          messageSender,
		chatCountWorkerPool:    chatCountWorkerPool,
		chatReadDateWorkerPool: chatReadDateWorkerPool,
	}
}

func (r *chatUsecase) RecvChatMessage(ctx context.Context, input input.ChatMessageInput) {

	chatRoomEntity := entity.MakeChatRoomEntity(input.ChatRoomData.RoomType, input.ChatRoomData.RoomKey, input.ChatRoomData.SecretFlag)
	chatLineEntity := entity.MakeChatLineEntity(input.ChatLineData.Cmd, input.ChatLineData.Contents, input.ChatLineData.LineKey, input.ChatLineData.TargetLineKey, input.ChatLineData.SendUserHash, input.ChatLineData.SendDate)
	chatFileEntity := make([]entity.ChatFileEntity, 0)
	for _, f := range input.ChatFileData {

		temp := entity.ChatFileEntity{
			FileId:       f.FileId,
			FileName:     f.FileName,
			FileType:     f.FileType,
			ThumbnailUrl: f.ThumbnailUrl,
		}
		chatFileEntity = append(chatFileEntity, temp)
	}

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

	chatFileOutput := make([]output.ChatFileDataOutput, 0)
	for _, f := range chatFileEntity {

		temp := output.ChatFileDataOutput{
			FileId:       f.FileId,
			FileName:     f.FileName,
			FileType:     f.FileType,
			ThumbnailUrl: f.ThumbnailUrl,
		}

		chatFileOutput = append(chatFileOutput, temp)
	}

	chatMessageOutput := output.ChatMessageOutput{
		ChatSession:  input.ChatSession,
		ChatRoomData: chatRoomOutput,
		ChatLineData: chatLintOutput,
		ChatFileData: chatFileOutput,
	}

	// 이 시점부터는 “이벤트 → 패킷” 변환
	out := dto.WsResponseDTO[output.ChatMessageOutput]{
		Type:      consts.CHAT,
		EventType: input.EventType,
		Data:      chatMessageOutput,
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

	chatCountMessageEntity := entity.ChatCountMessageEntity{
		Type:          consts.CHATUNREAD,
		EventType:     chatCountEntity.EventType,
		ChatCountData: chatCountData,
	}

	if chatCountMessageEntity.EventType == consts.READ {
		// 읽음처리 - 나에게 발송
		r.chatCountWorkerPool.AddTask(in.SendUserHash, chatCountMessageEntity)
	} else if chatCountMessageEntity.EventType == consts.UNREAD {
		// 신규 라인 발생 - 발신자를 제외하고 보냄.
		recvUserHash := r.chatRoomStorage.GetChatRoomMember(in.RoomKey)
		for _, recvUser := range recvUserHash {
			if recvUser != chatCountEntity.SendUserHash {
				r.chatCountWorkerPool.AddTask(recvUser, chatCountMessageEntity)
			}
		}
	}

}

func (r *chatUsecase) RecvChatReadMessage(ctx context.Context, in input.ChatReadMessageInput) {

	chatReadEntity := entity.MakeChatReadEntity(in.RoomKey, in.RoomType, in.MemberHash, in.ReadDate)
	log.Println("[RecvChatReadMessage] chatReadEntity: ", chatReadEntity)

	chatReadMessageEntity := entity.ChatReadMessageEntity{
		Type:         "readDate",
		ChatReadData: chatReadEntity,
	}
	recvUserHash := r.chatRoomStorage.GetChatRoomMember(chatReadMessageEntity.ChatReadData.RoomKey)

	// 나를 포함하여 보내야 중복 로그인 시 동일 사용자의 다른 클라이언트에 response 함.
	log.Println("1")
	for _, recvUser := range recvUserHash {
		r.chatReadDateWorkerPool.AddTask(recvUser, chatReadMessageEntity)
	}

}
