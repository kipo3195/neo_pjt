package chatRoom

type BulkCreateChatRoomDto struct {
	RoomKey     string `json:"roomKey"`
	RegDate     string `json:"regDate"`
	MemberCount int    `json:"memberCount"`
}

type BulkCreateChatRoomResponse struct {
	MakeCount    int                     `json:"makeCount"`
	CreatedCount int                     `json:"createdCount"`
	Room         []BulkCreateChatRoomDto `json:"room"`
}
