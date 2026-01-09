package entity

type ChatCountEntity struct {
	RoomKey      string
	RoomType     string
	EventType    string
	SendUserHash string
	Delta        int
}

func MakeChatCountEntity(roomKey string, roomType string, eventType string, sendUserHash string, delta int) ChatCountEntity {
	return ChatCountEntity{
		RoomKey:      roomKey,
		RoomType:     roomType,
		EventType:    eventType,
		SendUserHash: sendUserHash,
		Delta:        delta,
	}
}
