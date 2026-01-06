package repository

import (
	"context"
	"notificator/internal/domain/chatRoom/entity"
)

type ChatRoomRepository interface {
	GetMyChatRoom(userHash string) (entity.MyChatRoomEntity, error)
	PutChatRoomMember(ctx context.Context, en entity.CreateChatRoomEntity) error
}
