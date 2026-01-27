package input

type ChatMessageInput struct {
	EventType    string            `json:"eventType"`
	ChatSession  string            `json:"chatSession"`
	ChatRoomData ChatRoomDataInput `json:"chatRoomData"`
	ChatLineData ChatLineDataInput `json:"chatLineData"`
	ChatFileData []ChatFileInput   `json:"chatFileData"`
}
