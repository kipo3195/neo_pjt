package repository

import (
	"common/internal/domain/appToken/repository"

	"gorm.io/gorm"
)

type appValidationRepositoryImpl struct {
	db *gorm.DB
}

func NewAppValidationRepository(db *gorm.DB) repository.AppTokenRepository {
	return &appValidationRepositoryImpl{}
}
