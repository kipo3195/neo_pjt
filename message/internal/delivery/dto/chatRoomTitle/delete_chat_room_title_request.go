package chatRoomTitle

type DeleteChatRoomTitleRequest struct {
	Org     string `json:"org"`
	RoomKey string `json:"roomKey"`
	Type    string `json:"type"`
}
