package repository

import (
	"context"
	"notificator/internal/domain/chat/entity"
)

type ChatRepository interface {
	PutChatRoomMember(ctx context.Context, en entity.CreateChatRoomEntity) error
	GetMyChatRoom(userHash string) (entity.MyChatRoomEntity, error)
}
