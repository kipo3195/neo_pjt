package repository

import (
	"message/internal/domain/chatFile/repository"
	"message/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type chatFileRepositoryImpl struct {
	db *gorm.DB
}

func NewChatFileRepository(db *gorm.DB) repository.ChatFileRepository {
	return &chatFileRepositoryImpl{
		db: db,
	}
}

func ChatFileMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ChatFileHistory{})
}
