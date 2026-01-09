package chatService

type ReadChatRequest struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
}
