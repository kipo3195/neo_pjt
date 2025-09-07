package repository

import (
	"auth/internal/domain/userAuth/repository"

	"gorm.io/gorm"
)

type userAuthRepository struct {
	db *gorm.DB
}

func NewUserAuthRepository(db *gorm.DB) repository.UserAuthRepository {
	return &userAuthRepository{
		db: db,
	}
}
