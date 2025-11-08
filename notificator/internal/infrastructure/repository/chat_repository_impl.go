package repository

import (
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRepositoryImpl struct {
	db *gorm.DB
}

func ChatMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ChatMessage{})
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepositoryImpl{db: db}
}
