package repository

import (
	"context"
	"message/internal/domain/chatRoomTitle/entity"
)

type ChatRoomTitleRepository interface {
	UpdateChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error
	DeleteChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error
	GetChatRoomType(ctx context.Context, en entity.ChatRoomTitleEntity) (string, error)
}
