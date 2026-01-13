package entity

type GetChatRoomMemberReadDateEntity struct {
	RoomKey  string
	RoomType string
	UserHash string
}

func MakeGetChatRoomMemberReadDateEntity(roomKey string, roomType string, userHash string) GetChatRoomMemberReadDateEntity {
	return GetChatRoomMemberReadDateEntity{
		RoomKey:  roomKey,
		RoomType: roomType,
		UserHash: userHash,
	}
}
