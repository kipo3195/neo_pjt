package job

import "notificator/internal/domain/chat/entity"

type ChatReadDateJob struct {
	UserHash string
	En       entity.ChatReadMessageEntity
	IsFlush  bool
}
