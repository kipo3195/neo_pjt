package entity

type ChatRoomEntity struct {
	RoomType string
	RoomKey  string
}

func MakeChatRoomEntity(roomType string, roomKey string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType: roomType,
		RoomKey:  roomKey,
	}
}
