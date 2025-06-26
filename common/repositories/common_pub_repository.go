package repositories

import "gorm.io/gorm"

type commonPubRepository struct {
	db *gorm.DB
}

type CommonPubRepository interface {
}

func NewCommonPubRepository(db *gorm.DB) CommonPubRepository {
	return &commonPubRepository{db: db}
}
