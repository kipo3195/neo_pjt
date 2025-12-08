package chatRoom

type GetChatRoomListDto struct {
	ChatRoomDetail ChatRoomDetail `json:"roomDetail"`
	Member         []string       `json:"member"`
}
