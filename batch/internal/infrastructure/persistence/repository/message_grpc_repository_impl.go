package repository

import (
	"batch/internal/domain/messageGrpc/repository"

	"gorm.io/gorm"
)

type messageGrpcRepository struct {
	db *gorm.DB
}

func NewChatFileRepository(db *gorm.DB) repository.MessageGrpcRepository {
	return &messageGrpcRepository{
		db: db,
	}
}
