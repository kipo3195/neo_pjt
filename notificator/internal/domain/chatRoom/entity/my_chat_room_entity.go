package entity

type MyChatRoomEntity struct {
	RoomKey []string `gorm:"column:room_key"`
}
