package repository

import (
	"file/internal/domain/uploadFileCheck/repository"

	"gorm.io/gorm"
)

type uploadFileCheckRepositoryImpl struct {
	db *gorm.DB
}

func NewUploadFileCheckRepository(db *gorm.DB) repository.UploadFileCheckRepository {
	return &uploadFileCheckRepositoryImpl{
		db: db,
	}
}
