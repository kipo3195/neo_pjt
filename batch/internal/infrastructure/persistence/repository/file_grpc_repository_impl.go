package repository

import (
	"batch/internal/domain/fileGrpc/repository"

	"gorm.io/gorm"
)

type fileGrpcRepository struct {
	db *gorm.DB
}

func NewFileGrpcRepository(db *gorm.DB) repository.FileGrpcRepository {
	return &fileGrpcRepository{
		db: db,
	}
}
