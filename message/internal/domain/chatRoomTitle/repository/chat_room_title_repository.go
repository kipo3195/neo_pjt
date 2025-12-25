package repository

import (
	"context"
	"message/internal/domain/chatRoomTitle/entity"
)

type ChatRoomTitleRepository interface {
	DeleteChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error
	UpdateChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error
}
