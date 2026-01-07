package output

type ChatRoomMessageOutput struct {
	Type              string            `json:"type"`
	EventType         string            `json:"eventType"`
	ChatRoomEventData ChatRoomEventData `json:"chatRoomData"`
}
