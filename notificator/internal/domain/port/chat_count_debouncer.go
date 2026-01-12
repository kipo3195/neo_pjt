package port

import "notificator/internal/domain/chat/entity"

// 디바운싱은 네트워크 부하를 줄이기 위한 **전송 기술(Performance Strategy)**이다.
type ChatCountDebouncer interface {
	AddChatCount(userHash string, chatCountMessageEntity *entity.ChatCountMessageEntity) error
}
