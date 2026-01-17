package output

type GetChatRoomListOutput struct {
	ChatRoomDetail  ChatRoomDetail
	Member          []string
	MyChatRoomTitle ChatRoomTitleOutput
	Owner           ChatRoomOwnerOutput
	Line            ChatLineOutput
	Unread          ChatUnreadOutput
}
