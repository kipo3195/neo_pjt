package client

import "gorm.io/gorm"

type orgFileRepository struct {
	db *gorm.DB
}

type OrgFileRepository interface {
}

func NewOrgFileRepository(db *gorm.DB) OrgFileRepository {
	return &orgFileRepository{
		db: db,
	}
}
