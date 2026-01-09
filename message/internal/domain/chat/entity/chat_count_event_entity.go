package entity

type ChatCountEventEntity struct {
	RoomKey      string `json:"roomKey"`
	RoomType     string `json:"roomType"`
	EventType    string `json:"eventType"`
	SendUserHash string `json:"sendUserHash"`
	Delta        int    `json:"delta"` // 변경 건수 (EventType 따라 달라짐)
}

func MakeChatCountEventEntity(roomKey string, roomType string, eventType string, sendUserHash string, delta int) ChatCountEventEntity {

	return ChatCountEventEntity{
		RoomKey:      roomKey,
		RoomType:     roomType,
		EventType:    eventType,
		SendUserHash: sendUserHash,
		Delta:        delta,
	}
}
