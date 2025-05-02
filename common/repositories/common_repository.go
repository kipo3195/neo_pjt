package repositories

import (
	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

type CommonRepository interface {
}

func NewCommonRepository(db *gorm.DB) CommonRepository {

	return &commonRepository{db: db}
}
