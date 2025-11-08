package repository

import (
	"notificator/internal/domain/auth/repository"

	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepositoryImpl{db: db}
}
