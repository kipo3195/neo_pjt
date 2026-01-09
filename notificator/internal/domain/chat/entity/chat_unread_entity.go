package entity

type ChatUnreadEntity struct {
	RoomKey      string
	RoomType     string
	UnreadType   string
	SendUserHash string
	Delta        int
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
