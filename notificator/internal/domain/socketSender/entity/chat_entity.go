package entity

import "notificator/internal/consts"

type ChatEntity struct {
	Type           string
	EventType      string
	ChatSession    string
	ChatRoomEntity ChatRoomEntity
	ChatLineEntity ChatLineEntity
}

func MakeChatEntity(eventType string, chatSession string, chatRoomEntity ChatRoomEntity, chatLineEntity ChatLineEntity) ChatEntity {
	return ChatEntity{
		Type:           consts.CHAT,
		EventType:      eventType,
		ChatSession:    chatSession,
		ChatRoomEntity: chatRoomEntity,
		ChatLineEntity: chatLineEntity,
	}
}
