package chatRoomTitle

type UpdateChatRoomTitleRequest struct {
	Org     string `json:"org"`
	RoomKey string `json:"roomKey"`
	Type    string `json:"type"`
	Title   string `json:"title"`
}
