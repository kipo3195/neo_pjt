package repository

import (
	"context"
	"log"
	"notificator/internal/domain/chat/entity"
	"notificator/internal/domain/chat/repository"
	"notificator/internal/infrastructure/model"

	"gorm.io/gorm"
)

type chatRepositoryImpl struct {
	db *gorm.DB
}

func ChatMigrate(db *gorm.DB) {
	//db.AutoMigrate(&model.ChatMessage{})
	db.AutoMigrate(&model.ChatRoomMember{})
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepositoryImpl{db: db}
}

func (r *chatRepositoryImpl) PutChatRoomMember(ctx context.Context, en entity.CreateChatRoomEntity) error {

	memberModels := make([]model.ChatRoomMember, len(en.CreateChatRoomMember))
	for i, member := range en.CreateChatRoomMember {
		memberModels[i] = model.ChatRoomMember{
			RoomKey:    en.RoomKey,
			MemberHash: member.MemberHash,
		}
	}

	// bulk insert
	if err := r.db.
		WithContext(ctx).
		Create(&memberModels).Error; err != nil {
		return err
	}

	return nil
}

func (r *chatRepositoryImpl) GetMyChatRoom(userHash string) (entity.MyChatRoomEntity, error) {

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
