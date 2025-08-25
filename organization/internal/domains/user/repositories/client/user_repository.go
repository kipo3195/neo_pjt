package client

import "gorm.io/gorm"

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{
		db: db,
	}
}
