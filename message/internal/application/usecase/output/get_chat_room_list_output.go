package output

type GetChatRoomListOutput struct {
	ChatRoomDetail  ChatRoomDetail
	Member          []string
	MyChatRoomTitle ChatRoomTitleOutput
}
