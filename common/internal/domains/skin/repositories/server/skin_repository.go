package server

import "gorm.io/gorm"

type skinRepository struct {
	db *gorm.DB
}

type SkinRepository interface {
}

func NewSkinRepository(db *gorm.DB) SkinRepository {
	return &skinRepository{
		db: db,
	}
}
