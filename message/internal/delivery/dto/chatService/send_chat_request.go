package chatService

type SendChatRequest struct {
	EventType   string       `json:"eventType"`
	ChatLine    ChatLineData `json:"chatLineData"`
	ChatRoom    ChatRoomData `json:"chatRoomData"`
	ChatSession string       `json:"chatSession"`
}
