package repository

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/chatRoomTitle/entity"
	"message/internal/domain/chatRoomTitle/repository"
	"message/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRoomTitleRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomTitleRepository(db *gorm.DB) repository.ChatRoomTitleRepository {
	return &chatRoomTitleRepositoryImpl{
		db: db,
	}
}

func ChatRoomTitleMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ChatRoomTitle{})
}

func (r *chatRoomTitleRepositoryImpl) UpdateChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error {

	result := r.db.Exec(`
		INSERT INTO chat_room_title (org, user_hash, room_key, my_room_title, update_flag, update_date)
		SELECT 
			?, ?, ?, ?, 'Y', ?
		FROM 
			chat_room
		WHERE 
			room_key = ? AND room_type = ?
		ON DUPLICATE KEY UPDATE
			my_room_title = VALUES(my_room_title),
			update_flag   = 'Y',
			update_date   = VALUES(update_date)
		`,
		en.Org, en.UserHash, en.RoomKey, en.Title, en.EventDate,
		en.RoomKey, en.Type,
	)

	if result.Error != nil {
		log.Println("[UpdateChatRoomTitle] err : ", result.Error)
		return consts.ErrDB
	}

	if result.RowsAffected == 0 {
		log.Println("[UpdateChatRoomTitle] RowsAffected = 0")
		return consts.ErrDBresultNotFound
	}

	return nil
}

func (r *chatRoomTitleRepositoryImpl) DeleteChatRoomTitle(ctx context.Context, en entity.ChatRoomTitleEntity) error {

	result := r.db.
		Where("org = ? AND user_hash = ? AND room_key = ?",
			en.Org,
			en.UserHash,
			en.RoomKey,
		).
		Delete(&model.ChatRoomTitle{})

	if result.Error != nil {
		log.Println("[DeleteChatRoomTitle] err : ", result.Error)
		return consts.ErrDB
	}

	if result.RowsAffected == 0 {
		log.Println("[DeleteChatRoomTitle] RowsAffected = 0")
		return consts.ErrDBresultNotFound
	}

	return nil
}
