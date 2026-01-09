package output

type ChatCountDataOutput struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
	Delta    int    `json:"delta"`
}
