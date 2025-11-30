package entity

type ChatRoomEntity struct {
	RoomKey  string
	RoomType string
}

func MakeChatRoomEntity(roomKey string, roomType string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType: roomType,
		RoomKey:  roomKey,
	}
}
