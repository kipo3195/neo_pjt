package entity

type ChatRoomEntity struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	SecretFlag string `json:"secretFlag"`
}

func MakeChatRoomEntity(roomKey string, roomType string, secretFlag string) ChatRoomEntity {
	return ChatRoomEntity{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
