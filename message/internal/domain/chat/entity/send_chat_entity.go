package entity

type SendChatEntity struct {
	EventType      string           `json:"eventType"`
	ChatSession    string           `json:"chatSession"`
	ChatRoomEntity ChatRoomEntity   `json:"chatRoomData"`
	ChatLineEntity ChatLineEntity   `json:"chatLineData"`
	ChatFileEntity []ChatFileEntity `json:"chatFileData,omitempty"`
}

func MakeSendChatEntity(eventType string, chatSession string, chatLineEntity ChatLineEntity, chatRoomEntity ChatRoomEntity, chatFileEntity []ChatFileEntity) SendChatEntity {

	return SendChatEntity{
		EventType:      eventType,
		ChatSession:    chatSession,
		ChatLineEntity: chatLineEntity,
		ChatRoomEntity: chatRoomEntity,
		ChatFileEntity: chatFileEntity,
	}
}
