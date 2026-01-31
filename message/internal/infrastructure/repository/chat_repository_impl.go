package repository

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/chat/cache"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRepository struct {
	db           *gorm.DB
	cacheStorage cache.ChatCache
}

func NewChatRepository(db *gorm.DB, cacheStorage cache.ChatCache) repository.ChatRepository {
	return &chatRepository{
		db:           db,
		cacheStorage: cacheStorage,
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

	if len(sendChatEntity.ChatFileEntity) > 0 {

		// 채팅 파일
		chatFileHistory := make([]model.ChatFileHistory, len(sendChatEntity.ChatFileEntity))
		for i, file := range sendChatEntity.ChatFileEntity {

			var fileType string
			if file.FileType == consts.IMAGE {
				fileType = "img"
			} else {
				fileType = "file"
			}

			chatFileHistory[i] = model.ChatFileHistory{
				RoomKey:     sendChatEntity.ChatRoomEntity.RoomKey,
				LineKey:     sendChatEntity.ChatLineEntity.LineKey,
				FileId:      file.FileId,
				FileName:    file.FileName,
				FileType:    fileType,
				ReqUserHash: sendChatEntity.ChatLineEntity.SendUserHash,
			}
		}

		if err := tx.Create(chatFileHistory).Error; err != nil {
			tx.Rollback()
			log.Println("[SaveChatLine] unread update failed:", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Println("[SaveChatLine] - DB process Success lineKey : ", sendChatEntity.ChatLineEntity.LineKey)
	return nil
}

func (r *chatRepository) ReadChatLine(ctx context.Context, readChatEntity entity.ReadChatEntity) error {

	result := r.db.WithContext(ctx).Model(&model.ChatRoomMember{}).
		Where("room_key = ?", readChatEntity.RoomKey).
		Where("member_hash = ?", readChatEntity.UserHash).
		Where("member_state = ?", "1").
		Where("member_unread_count > 0"). // 거의 동시 타이밍에 읽음처리 됬을때 뒤에 들어온 요청을 무시(notifcator로 전달 X) 하기 위함.
		Updates(map[string]interface{}{
			// SQL의 컬럼 연산을 위해 gorm.Expr 사용
			"member_unread_count":      0,
			"member_unread_count_date": readChatEntity.ReadDate,
			"member_read_date":         readChatEntity.ReadDate,
		})

	if result.Error != nil {
		log.Printf("[ReadChatLine] - %s member unread update failed :%s\n", readChatEntity.UserHash, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("[ReadChatLine] - %s member unread count not update \n", readChatEntity.UserHash)
		return consts.ErrDBResultNotUpdate
	}

	log.Println("[ReadChatLine] - DB process Success userHash: ", readChatEntity.UserHash)
	return nil
}

func (r *chatRepository) GetChatLineEvent(ctx context.Context, en entity.GetChatLineEventEntity) ([]entity.ChatLineEventEntity, error) {

	var result []entity.ChatLineEventEntity

	err := r.db.Raw(
		`select 
			event_type, cmd, event.line_key, target_line_key, contents, send_user_hash, event.send_date,
			file.file_id, file.file_name, file.file_type
		from 
			chat_line_event as event 
		join 
			(select room.room_key
			from chat_room as room join chat_room_member as member 
			on room.room_key = member.room_key and member.member_hash = ?
			where room.room_key = ? ) as room_view 
		on
			event.room_key = room_view.room_key
		left join 
			(select room_key, line_key, file_id, file_name, file_type
			from chat_file_history 
			where room_key = ? ) as file
		on 
			event.line_key = file.line_key
		where 
			event.line_key > ? order by send_date asc`, en.ReqUserHash, en.RoomKey, en.RoomKey, en.LineKey).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *chatRepository) GetChatFileEntity(ctx context.Context, transactionId string) ([]*entity.ChatFileEntity, error) {

	return r.cacheStorage.GetFileEntity(ctx, transactionId)
}
