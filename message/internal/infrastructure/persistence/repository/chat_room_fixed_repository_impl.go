package repository

import (
	"message/internal/domain/chatRoomFixed/repository"

	"gorm.io/gorm"
)

type chatRoomFixedRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomFixedRepository(db *gorm.DB) repository.ChatRoomFixedRepository {
	return &chatRoomFixedRepositoryImpl{
		db: db,
	}
}
