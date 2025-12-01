package entity

type ChatRoomEntity struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	SecretFlag bool   `json:"secretFlag"`
}

func MakeChatRoomEntity(roomKey string, roomType string, secretFlag bool) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
