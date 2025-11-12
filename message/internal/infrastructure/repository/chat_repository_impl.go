package repository

import (
	"message/internal/domain/chat/repository"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepository{
		db: db,
	}
}
