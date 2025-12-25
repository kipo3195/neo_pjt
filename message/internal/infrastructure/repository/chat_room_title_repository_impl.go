package repository

import (
	"context"
	"message/internal/domain/chatRoomTitle/entity"
	"message/internal/domain/chatRoomTitle/repository"

	"gorm.io/gorm"
)

type chatRoomTitleRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomTitleRepository(db *gorm.DB) repository.ChatRoomTitleRepository {
	return &chatRoomTitleRepositoryImpl{
		db: db,
	}
}

func (r *chatRoomTitleRepositoryImpl) UpdateChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error {

	return nil
}

func (r *chatRoomTitleRepositoryImpl) DeleteChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error {

	return nil
}
