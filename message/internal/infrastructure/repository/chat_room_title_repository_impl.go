package repository

import (
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
