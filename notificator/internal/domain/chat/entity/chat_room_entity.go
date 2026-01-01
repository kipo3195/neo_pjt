package entity

type ChatRoomEntity struct {
	RoomType   string
	RoomKey    string
	SecretFlag string
}

func MakeChatRoomEntity(roomType string, roomKey string, secretFlag string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
