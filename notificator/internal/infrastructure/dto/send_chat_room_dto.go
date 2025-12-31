package dto

type SendChatRoomDto struct {
	RoomType   string `json:"roomType"`
	RoomKey    string `json:"roomKey"`
	SecretFlag bool   `json:"secretFlag"`
}

func MakeChatRoomDto(roomType string, roomKey string, secretFlag bool) SendChatRoomDto {
	return SendChatRoomDto{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
