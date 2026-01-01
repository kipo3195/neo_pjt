package dto

type SendChatRoomDto struct {
	RoomType   string `json:"roomType"`
	RoomKey    string `json:"roomKey"`
	SecretFlag string `json:"secretFlag"`
}

func MakeChatRoomDto(roomType string, roomKey string, secretFlag string) SendChatRoomDto {
	return SendChatRoomDto{
		RoomType:   roomType,
		RoomKey:    roomKey,
		SecretFlag: secretFlag,
	}
}
