package job

import "notificator/internal/domain/chat/entity"

type ChatCountJob struct {
	UserHash string
	En       entity.ChatCountMessageEntity
	IsFlush  bool // 타이머에 의해 발생한 플러시 작업인지 구분
}
