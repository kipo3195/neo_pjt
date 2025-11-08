package repository

import "gorm.io/gorm"

type chatRepository struct {
	db *gorm.DB
}

type ChatRepository interface {
}
