package dto

type ChatCountDataDto struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
	Delta    int    `json:"delta"`
}
