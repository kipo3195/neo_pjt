package repository

import (
	"file/internal/domain/chatFile/repository"

	"gorm.io/gorm"
)

type chatFileRespositoryImpl struct {
	db *gorm.DB
}

func NewChatFileRepository(db *gorm.DB) repository.ChatFileRepository {
	return &chatFileRespositoryImpl{
		db: db,
	}
}
