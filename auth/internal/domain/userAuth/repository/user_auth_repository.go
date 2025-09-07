package repository

import "gorm.io/gorm"

type userAuthRepository struct {
	db *gorm.DB
}

type UserAuthRepository interface {
}
