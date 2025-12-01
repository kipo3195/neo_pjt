package entity

import "notificator/internal/consts"

type SendChatEntity struct {
	Type               string
	EventType          string
	ChatSession        string
	SendChatRoomEntity SendChatRoomEntity
	SendChatLineEntity SendChatLineEntity
}

func MakeSendChatEntity(eventType string, chatSession string, sendChatRoomEntity SendChatRoomEntity, sendChatLineEntity SendChatLineEntity) SendChatEntity {
	return SendChatEntity{
		Type:               consts.CHAT,
		EventType:          eventType,
		ChatSession:        chatSession,
		SendChatRoomEntity: sendChatRoomEntity,
		SendChatLineEntity: sendChatLineEntity,
	}
}
