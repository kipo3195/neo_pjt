package entity

type ChatReadEntity struct {
	RoomKey    string
	RoomType   string
	MemberHash string
	ReadDate   string
}

func MakeChatReadEntity(roomKey string, roomType string, memberHash string, readDate string) ChatReadEntity {
	return ChatReadEntity{
		RoomKey:    roomKey,
		RoomType:   roomType,
		MemberHash: memberHash,
		ReadDate:   readDate,
	}
}
