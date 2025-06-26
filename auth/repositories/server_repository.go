package repositories

import (
	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

type ServerRepository interface {
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{
		db: db,
	}
}
