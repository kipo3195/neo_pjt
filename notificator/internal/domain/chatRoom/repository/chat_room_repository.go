package repository

import "notificator/internal/domain/chatRoom/entity"

type ChatRoomRepository interface {
	GetMyChatRoom(userHash string) (entity.MyChatRoomEntity, error)
}
