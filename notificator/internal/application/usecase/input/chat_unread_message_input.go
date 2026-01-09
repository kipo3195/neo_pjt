package input

type ChatUnreadMessageInput struct {
	RoomKey      string `json:"roomKey"`
	RoomType     string `json:"roomType"`
	UnreadType   string `json:"unreadType"`
	SendUserHash string `json:"sendUserHash"`
	Delta        int    `json:"delta"` // 변경 건수 (unreadType에 따라 달라짐)
}
