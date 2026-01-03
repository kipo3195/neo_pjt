package entity

type ChatRoomUpdateDateEntity struct {
	RoomKey  string `gorm:"column:room_key"`
	RoomType string `gorm:"column:room_type"`
	Detail   string `gorm:"column:detail"`
	Line     string `gorm:"column:line"`
	Member   string `gorm:"column:member"`
	Owner    string `gorm:"column:owner"`
	Title    string `gorm:"column:title"`
}
