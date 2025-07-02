package repositories

import "gorm.io/gorm"

type publicRepository struct {
	db *gorm.DB
}

type PublicRepository interface {
}

func NewPublicRepository(db *gorm.DB) PublicRepository {
	return &publicRepository{db: db}
}
