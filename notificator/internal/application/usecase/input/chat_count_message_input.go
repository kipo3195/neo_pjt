package input

type ChatCountMessageInput struct {
	RoomKey      string `json:"roomKey"`
	RoomType     string `json:"roomType"`
	EventType    string `json:"eventType"`
	SendUserHash string `json:"sendUserHash"`
	Delta        int    `json:"delta"` // 변경 건수 (unreadType에 따라 달라짐)
}
