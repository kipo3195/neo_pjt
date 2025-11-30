package entity

type ChatRoomEntity struct {
	RoomKey  string `json:"roomKey"`
	RoomType string `json:"roomType"`
}

func MakeChatRoomEntity(roomKey string, roomType string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType: roomType,
		RoomKey:  roomKey,
	}
}
