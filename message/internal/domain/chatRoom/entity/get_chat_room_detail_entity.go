package entity

type GetChatRoomDetailEntity struct {
	ReqUserHash string
	RoomType    string
	RoomKey     []string
}

func MakeGetChatRoomDetailEntity(reqUserHash string, roomType string, roomKey []string) GetChatRoomDetailEntity {
	return GetChatRoomDetailEntity{
		ReqUserHash: reqUserHash,
		RoomType:    roomType,
		RoomKey:     roomKey,
	}
}
