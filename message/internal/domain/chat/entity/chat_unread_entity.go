package entity

type ChatUnreadEntity struct {
	RoomKey      string `json:"roomKey"`
	RoomType     string `json:"roomType"`
	UnreadType   string `json:"unreadType"`
	SendUserHash string `json:"sendUserHash"`
	Delta        int    `json:"delta"` // 변경 건수 (unreadType에 따라 달라짐)
}

func MakeChatUnreadEntity(roomKey string, roomType string, unreadType string, sendUserHash string, delta int) ChatUnreadEntity {

	return ChatUnreadEntity{
		RoomKey:      roomKey,
		RoomType:     roomType,
		UnreadType:   unreadType,
		SendUserHash: sendUserHash,
		Delta:        delta,
	}
}
