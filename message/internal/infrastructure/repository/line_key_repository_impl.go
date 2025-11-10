package repository

import (
	"message/internal/domain/lineKey/repository"

	"gorm.io/gorm"
)

type lineKeyRepository struct {
	db *gorm.DB
}

func NewLineKeyRepository(db *gorm.DB) repository.LineKeyRepository {
	return &lineKeyRepository{
		db: db,
	}
}
