package chatRoom

type GetChatRoomDetailRequest struct {
	RoomType string   `json:"roomType" validate:"required"`
	RoomKey  []string `json:"roomKey" validate:"required"`
}
