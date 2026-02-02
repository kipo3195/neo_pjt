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

	// db.Transaction의 핵심은 클로저 내부에서 error가 반환되면 GORM이 알아서 롤백을 수행한다는 점입니다.
	// 에러가 발생했을 때 그냥 return err만 해주시면 됩니다.
	// WithContext(ctx)를 먼저 호출하고 Transaction을 실행하면, 내부의 tx는 이미 해당 ctx를 가진 상태입니다.
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// ❌ 틀림: r.db는 트랜잭션 외부에 있는 객체입니다.
		// r.db.WithContext(ctx).Create(&data) <- 이렇게 쓰면 트랜잭션에 포함되지 않습니다. 반드시 클로저의 인자로 들어온 tx를 사용해야만, 해당 tx가 관리하는 동일한 데이터베이스 커넥션과 컨텍스트 안에서 모든 작업이 안전하게 묶이게 됩니다.

		// 채팅 라인 저장
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
			log.Println("[SaveChatLine] line insert failed:", err)
			return err
		}

		// 읽지 않은 메시지 수 업데이트
		if err := tx.Model(&model.ChatRoomMember{}).
			Where("room_key = ?", sendChatEntity.ChatRoomEntity.RoomKey).
			Where("member_hash != ?", sendChatEntity.ChatLineEntity.SendUserHash).
			Where("member_state = ?", "1").
			Updates(map[string]interface{}{
				"member_unread_count":      gorm.Expr("member_unread_count + 1"),
				"member_unread_count_date": sendChatEntity.ChatLineEntity.SendDate,
			}).Error; err != nil {
			log.Println("[SaveChatLine] unread update failed:", err)
			return err
		}

		// 채팅 파일 히스토리 저장
		if len(sendChatEntity.ChatFileEntity) > 0 {

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
				log.Println("[SaveChatLine] unread update failed:", err)
				return err
			}
		}
		return nil // nil을 반환하면 자동 Commit 됩니다.
	})

	// 트랜잭션 내부에서 발생한 에러가 여기까지 전달됩니다.
	if err != nil {
		log.Println(err)
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
