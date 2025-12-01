package input

type ChatRoomDataInput struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	SecretFlag bool   `json:"secretFlag"`
}
