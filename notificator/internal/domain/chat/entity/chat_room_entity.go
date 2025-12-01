package entity

type ChatRoomEntity struct {
	RoomType   string
	RoomKey    string
	SecretFlag bool
}

func MakeChatRoomEntity(roomType string, roomKey string, secretFlag bool) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
