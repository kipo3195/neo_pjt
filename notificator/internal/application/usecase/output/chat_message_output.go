package output

type ChatMessageOutput struct {
	EventType    string
	ChatSession  string
	ChatRoomData ChatRoomDataOutput
	ChatLineData ChatLineDataOutput
}
