package repository

import (
	"common/internal/domain/appToken/repository"

	"gorm.io/gorm"
)

type appTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewAppTokenRepository(db *gorm.DB) repository.AppTokenRepository {
	return &appTokenRepositoryImpl{}
}
