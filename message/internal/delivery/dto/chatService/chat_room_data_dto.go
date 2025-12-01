package chatService

type ChatRoomData struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	SecretFlag bool   `json:"secretFlag"`
}
