package client

import "gorm.io/gorm"

type appTokenRepository struct {
	db *gorm.DB
}

type AppTokenRepository interface {
}

func NewAppTokenRepository(db *gorm.DB) AppTokenRepository {
	return appTokenRepository{
		db: db,
	}
}
