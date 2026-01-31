package input

type SendChatInput struct {
	ChatRoom      ChatRoomInput
	ChatLine      ChatLineInput
	EventType     string
	ChatSession   string
	TransactionId string
}
