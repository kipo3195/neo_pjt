package repository

import (
	"file/internal/domain/fileUrl/repository"

	"gorm.io/gorm"
)

type fileUrlRepositoryImpl struct {
	db *gorm.DB
}

func NewFileUrlRepository(db *gorm.DB) repository.FileUrlRepository {

	return &fileUrlRepositoryImpl{
		db: db,
	}
}
