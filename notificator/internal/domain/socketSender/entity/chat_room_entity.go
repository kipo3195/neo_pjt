package entity

type ChatRoomEntity struct {
	RoomKey    string
	RoomType   string
	SecretFlag bool
}

func MakeChatRoomEntity(roomKey string, roomType string, secretFlag bool) ChatRoomEntity {
	return ChatRoomEntity{
		RoomKey:    roomKey,
		RoomType:   roomType,
		SecretFlag: secretFlag,
	}
}
