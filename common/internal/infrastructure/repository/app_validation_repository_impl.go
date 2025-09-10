package repository

import (
	"common/internal/domain/appToken/repository"

	"gorm.io/gorm"
)

type appValidationRepositoryImpl struct {
	db *gorm.DB
}

func NewAppTokenRepository(db *gorm.DB) repository.AppTokenRepository {
	return &appValidationRepositoryImpl{}
}
