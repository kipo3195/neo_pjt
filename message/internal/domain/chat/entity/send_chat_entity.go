package entity

type SendChatEntity struct {
	EventType      string         `json:"eventType"`
	ChatSession    string         `json:"chatSession"`
	ChatRoomEntity ChatRoomEntity `json:"chatRoomData"`
	ChatLineEntity ChatLineEntity `json:"chatLineData"`
}

func MakeSendChatEntity(eventType string, chatSession string, chatLineEntity ChatLineEntity, chatRoomEntity ChatRoomEntity) SendChatEntity {

	return SendChatEntity{
		EventType:      eventType,
		ChatSession:    chatSession,
		ChatLineEntity: chatLineEntity,
		ChatRoomEntity: chatRoomEntity,
	}
}
