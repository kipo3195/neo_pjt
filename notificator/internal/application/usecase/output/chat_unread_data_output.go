package output

type ChatUnreadDataOutput struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
	Delta    int    `json:"delta"`
}
