package repository

import (
	"message/internal/domain/chatRoomConfig/repository"

	"gorm.io/gorm"
)

type chatRoomConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomConfigRepository(db *gorm.DB) repository.ChatRoomConfigRepository {
	return &chatRoomConfigRepositoryImpl{
		db: db,
	}
}
