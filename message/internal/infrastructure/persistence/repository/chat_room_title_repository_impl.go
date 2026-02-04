package repository

import (
	"context"
	"errors"
	"log"
	"message/internal/consts"
	"message/internal/domain/chatRoomTitle/entity"
	"message/internal/domain/chatRoomTitle/repository"
	"message/internal/infrastructure/persistence/model"

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

	// 생성된 방을 기준으로 처리하기 위함 (insert select)
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

func (r *chatRoomTitleRepositoryImpl) GetChatRoomType(ctx context.Context, en entity.ChatRoomTitleEntity) (string, error) {

	var roomType string

	err := r.db.WithContext(ctx).
		Table("chat_room as room").
		Select("room.room_type").
		Joins("join chat_room_member as member on room.room_key = member.room_key and member.member_hash = ?", en.UserHash).
		Where("room.room_key = ?", en.RoomKey).
		Take(&roomType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", consts.ErrDBresultNotFound
		}
		return "", consts.ErrDB
	}

	return roomType, nil

}
