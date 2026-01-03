package repository

import (
	"context"
	"message/internal/domain/chatRoom/entity"
)

type ChatRoomRepository interface {
	PutChatRoom(ctx context.Context, entity entity.CreateChatRoomEntity) error
	GetChatRoomDetail(ctx context.Context, entity entity.GetChatRoomDetailEntity) ([]entity.ChatRoomDetailEntity, error)
	GetChatRoomList(ctx context.Context, entity entity.GetChatRoomListEntity) ([]entity.ChatRoomDetailEntity, error)
	GetChatRoomUpdateDate(ctx context.Context, entity entity.GetChatRoomUpdateDateEntity) ([]entity.ChatRoomUpdateDateEntity, error)
}
