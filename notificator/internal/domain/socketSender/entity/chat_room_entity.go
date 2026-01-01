package entity

type ChatRoomEntity struct {
	RoomKey    string
	RoomType   string
	SecretFlag string
}

func MakeChatRoomEntity(roomKey string, roomType string, secretFlag string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomKey:    roomKey,
		RoomType:   roomType,
		SecretFlag: secretFlag,
	}
}
