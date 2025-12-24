package chatRoom

type CreateChatRoomResponse struct {
	RoomKey string `json:"roomKey"`
	RegDate string `json:"regDate"`
	Type    string `json:"roomType"`
}
