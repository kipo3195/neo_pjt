package chatRoom

type GetChatRoomDetailDto struct {
	ChatRoomDetail ChatRoomDetail `json:"roomDetail"`
	Member         []string       `json:"member"`
}
