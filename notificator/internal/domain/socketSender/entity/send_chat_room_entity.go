package entity

type SendChatRoomEntity struct {
	RoomKey    string
	RoomType   string
	SecretFlag bool
}

func MakeSendChatRoomEntity(roomKey string, roomType string, secretFlag bool) SendChatRoomEntity {
	return SendChatRoomEntity{
		RoomKey:    roomKey,
		RoomType:   roomType,
		SecretFlag: secretFlag,
	}
}
