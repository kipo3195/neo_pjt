package entity

type SendChatEntity struct {
	EventType      string
	ChatSession    string
	ChatLineEntity ChatLineEntity
	ChatRoomEntity ChatRoomEntity
}

func MakeSendChatEntity(eventType string, chatSession string, chatLineEntity ChatLineEntity, chatRoomEntity ChatRoomEntity) SendChatEntity {

	return SendChatEntity{
		EventType:      eventType,
		ChatSession:    chatSession,
		ChatLineEntity: chatLineEntity,
		ChatRoomEntity: chatRoomEntity,
	}
}
