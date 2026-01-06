package repository

import (
	"log"
	"notificator/internal/domain/chatRoom/entity"
	"notificator/internal/domain/chatRoom/repository"
	"notificator/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRoomRepositoryImpl struct {
	db *gorm.DB
}

func ChatRoomMigrate(db *gorm.DB) {
	//db.AutoMigrate(&model.ChatMessage{})
	db.AutoMigrate(&model.ChatRoomMember{})
}

func NewChatRoomRepository(db *gorm.DB) repository.ChatRoomRepository {
	return &chatRoomRepositoryImpl{
		db: db,
	}
}

func (r *chatRoomRepositoryImpl) GetMyChatRoom(userHash string) (entity.MyChatRoomEntity, error) {

	var roomKey []string

	err := r.db.Raw(
		`select 
			room_key
		from chat_room_member
		where 
			member_hash = ? and member_state = '1' `,
		userHash).Scan(&roomKey).Error

	if err != nil {
		log.Println("[GetMyChatRoom] db error")
		return entity.MyChatRoomEntity{}, err
	}

	result := entity.MyChatRoomEntity{
		RoomKey: roomKey,
	}

	return result, nil
}
