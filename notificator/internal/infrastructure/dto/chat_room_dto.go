package dto

type ChatRoomDto struct {
	RoomType   string `json:"roomType"`
	RoomKey    string `json:"roomKey"`
	SecretFlag bool   `json:"secretFlag"`
}

func MakeChatRoomDto(roomType string, roomKey string, secretFlag bool) ChatRoomDto {
	return ChatRoomDto{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
