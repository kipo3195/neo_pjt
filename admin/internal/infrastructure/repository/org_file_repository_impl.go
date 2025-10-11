package repository

import (
	"admin/internal/domain/orgFile/repository"

	"gorm.io/gorm"
)

type orgFileRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgFileRepository(db *gorm.DB) repository.OrgFileRepository {
	return &orgFileRepositoryImpl{
		db: db,
	}
}
