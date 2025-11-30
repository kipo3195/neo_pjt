package repository

import (
	"message/internal/domain/chatRoom/repository"

	"gorm.io/gorm"
)

type chatRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) repository.ChatRoomRepository {
	return &chatRoomRepositoryImpl{
		db: db,
	}
}
