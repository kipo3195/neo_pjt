package repository

import (
	"context"
	"log"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func ChatLineEventMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ChatLineEvent{})
}

func (r *chatRepository) SaveChatLine(ctx context.Context, sendChatEntity entity.SendChatEntity) {

	err := r.db.WithContext(ctx).Create(&model.ChatLineEvent{
		EventType:    sendChatEntity.EventType,
		Cmd:          sendChatEntity.ChatLineEntity.Cmd,
		LineKey:      sendChatEntity.ChatLineEntity.LineKey,
		Contents:     sendChatEntity.ChatLineEntity.Contents,
		SendUserHash: sendChatEntity.ChatLineEntity.SendUserHash,
		SendDate:     sendChatEntity.ChatLineEntity.SendDate,
	}).Error

	if err != nil {
		log.Println("[SaveChatLine] - DB insert failed :", err)
	} else {
		log.Println("[SaveChatLine] - Insert Success")
	}

}
