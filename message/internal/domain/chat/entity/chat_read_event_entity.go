package entity

type ChatReadEventEntity struct {
	RoomKey    string `json:"roomKey"`
	RoomType   string `json:"roomType"`
	MemberHash string `json:"memberHash"`
	ReadDate   string `json:"readDate"`
}

func MakeChatReadEventEntity(roomKey string, roomType string, memberHash string, readDate string) ChatReadEventEntity {
	return ChatReadEventEntity{
		RoomKey:    roomKey,
		RoomType:   roomType,
		MemberHash: memberHash,
		ReadDate:   readDate,
	}
}
