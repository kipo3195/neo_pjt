package entity

type ChatRoomTitleEntity struct {
	UserHash  string
	RoomKey   string
	Org       string
	Type      string
	Title     string
	EventDate string // 처리시간
}

func MakeDeleteChatRoomTitleEntity(userhash string, org string, roomKey string, t string) ChatRoomTitleEntity {
	return ChatRoomTitleEntity{
		UserHash: userhash,
		Org:      org,
		RoomKey:  roomKey,
		Type:     t,
	}
}

func MakeUpdateChatRoomTitleEntity(userhash string, org string, roomKey string, t string, title string) ChatRoomTitleEntity {
	return ChatRoomTitleEntity{
		UserHash: userhash,
		Org:      org,
		RoomKey:  roomKey,
		Type:     t,
		Title:    title,
	}
}
