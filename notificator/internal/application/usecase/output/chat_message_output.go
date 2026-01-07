package output

type ChatMessageOutput struct {
	Type         string             `json:"type"`
	EventType    string             `json:"eventType"`
	ChatSession  string             `json:"chatSession"`
	ChatRoomData ChatRoomDataOutput `json:"chatRoomData"`
	ChatLineData ChatLineDataOutput `json:"chatLineData"`
}
