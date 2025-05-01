package repositories

import "gorm.io/gorm"

type chatRepository struct {
	db *gorm.DB
}

type ChatRepository interface {
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}
