package output

type BulkCreateChatRoomItemOutput struct {
	RoomKey     string
	RegDate     string
	MemberCount int
}

type BulkCreateChatRoomOutput struct {
	MakeCount    int
	CreatedCount int
	Room         []BulkCreateChatRoomItemOutput
}
