package output

type ChatMessageOutput struct {
	ChatSession  string             `json:"chatSession"`
	ChatRoomData ChatRoomDataOutput `json:"chatRoomData"`
	ChatLineData ChatLineDataOutput `json:"chatLineData"`
}
