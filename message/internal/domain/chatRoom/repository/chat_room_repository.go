package repository

import (
	"context"
	"message/internal/domain/chatRoom/entity"
)

type ChatRoomRepository interface {
	PutChatRoom(ctx context.Context, memberEntity []entity.CreateChatRoomMemberEntity, roomEntity entity.ChatRoomEntity) error
	GetChatRoomDetail(ctx context.Context, entity entity.GetChatRoomDetailEntity) ([]entity.ChatRoomDetailEntity, error)
	GetChatRoomList(ctx context.Context, entity entity.GetChatRoomListEntity) ([]entity.ChatRoomDetailEntity, error)
}
