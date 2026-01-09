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

func (r *chatRepository) SaveChatLine(ctx context.Context, sendChatEntity entity.SendChatEntity) error {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 실패 시 롤백 보장
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&model.ChatLineEvent{
		EventType:     sendChatEntity.EventType,
		Cmd:           sendChatEntity.ChatLineEntity.Cmd,
		RoomKey:       sendChatEntity.ChatRoomEntity.RoomKey,
		TargetLineKey: sendChatEntity.ChatLineEntity.TargetLineKey,
		LineKey:       sendChatEntity.ChatLineEntity.LineKey,
		Contents:      sendChatEntity.ChatLineEntity.Contents,
		SendUserHash:  sendChatEntity.ChatLineEntity.SendUserHash,
		SendDate:      sendChatEntity.ChatLineEntity.SendDate,
	}).Error; err != nil {
		tx.Rollback()
		log.Println("[SaveChatLine] line insert failed:", err)
		return err
	}

	if err := tx.Model(&model.ChatRoomMember{}).
		Where("room_key = ?", sendChatEntity.ChatRoomEntity.RoomKey).
		Where("member_hash != ?", sendChatEntity.ChatLineEntity.SendUserHash).
		Where("member_state = ?", "1").
		Updates(map[string]interface{}{
			"member_unread_count":      gorm.Expr("member_unread_count + 1"),
			"member_unread_count_date": sendChatEntity.ChatLineEntity.SendDate,
		}).Error; err != nil {
		tx.Rollback()
		log.Println("[SaveChatLine] unread update failed:", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Println("[SaveChatLine] - DB process Success lineKey : ", sendChatEntity.ChatLineEntity.LineKey)
	return nil
}

func (r *chatRepository) ReadChatLine(ctx context.Context, readChatEntity entity.ReadChatEntity) error {

	err := r.db.WithContext(ctx).Model(&model.ChatRoomMember{}).
		Where("room_key = ?", readChatEntity.RoomKey).
		Where("member_hash = ?", readChatEntity.UserHash).
		Where("member_state = ?", "1").
		Updates(map[string]interface{}{
			// SQL의 컬럼 연산을 위해 gorm.Expr 사용
			"member_unread_count":      0,
			"member_unread_count_date": readChatEntity.ReadDate,
			"member_read_date":         readChatEntity.ReadDate,
		}).Error

	if err != nil {
		log.Println("[ReadChatLine] - member unread update failed :", err)
		return err
	} else {

		log.Println("[ReadChatLine] - DB process Success userHash: ", readChatEntity.UserHash)
		return nil
	}

}
