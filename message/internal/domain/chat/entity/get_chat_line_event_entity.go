package entity

type GetChatLineEventEntity struct {
	ReqUserHash string
	Org         string
	RoomType    string
	RoomKey     string
	LineKey     string
}

func MakeGetChatLineEventEntity(reqUserHash string, org string, roomType string, roomKey string, lineKey string) GetChatLineEventEntity {
	return GetChatLineEventEntity{
		ReqUserHash: reqUserHash,
		Org:         org,
		RoomType:    roomType,
		RoomKey:     roomKey,
		LineKey:     lineKey,
	}
}
