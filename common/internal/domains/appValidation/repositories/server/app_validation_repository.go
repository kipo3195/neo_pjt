package server

import "gorm.io/gorm"

type appValidationRepository struct {
	db *gorm.DB
}

type AppValidationRepository interface {
}

func NewAppValidationRepository(db *gorm.DB) AppValidationRepository {
	return &appValidationRepository{
		db: db,
	}
}
