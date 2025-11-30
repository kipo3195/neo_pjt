package chatService

type SendChatResponse struct {
	ChatRoom    ChatRoomData `json:"chatRoomData"`
	ChatLine    ChatLineData `json:"chatLineData"`
	EventType   string       `json:"eventType"`
	ChatSession string       `json:"chatSession"`
}
