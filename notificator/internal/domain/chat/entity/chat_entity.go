package entity

type ChatMessage struct {
	EventType    string         `json:"eventType"`
	ChatSession  string         `json:"chatSession"`
	ChatRoomData ChatRoomEntity `json:"chatRoomData"`
	ChatLineData ChatLineEntity `json:"chatLineData"`
}

func MakeRecvChatMessageEntity(eventType string, chatSession string, chatRoomData ChatRoomEntity, chatLineData ChatLineEntity) ChatMessage {

	return ChatMessage{
		EventType:    eventType,
		ChatSession:  chatSession,
		ChatRoomData: chatRoomData,
		ChatLineData: chatLineData,
	}

}
