package repository

import (
	"notificator/internal/domain/login/repository"

	"gorm.io/gorm"
)

type loginRepositoryImpl struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) repository.LoginRepository {

	return &loginRepositoryImpl{
		db: db,
	}
}
