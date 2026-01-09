package entity

type ReadChatEntity struct {
	RoomKey  string
	RoomType string
	UserHash string
	ReadDate string
}

func MakeReadChatEntity(roomKey string, roomType string, userHash string, readDate string) ReadChatEntity {
	return ReadChatEntity{
		RoomKey:  roomKey,
		RoomType: roomType,
		UserHash: userHash,
		ReadDate: readDate,
	}
}
