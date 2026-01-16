package entity

import "time"

type ChatRoomMyEntity struct {
	RoomKey     string     `gorm:"column:room_key"`
	RoomType    string     `gorm:"column:room_type"`
	ReadDate    *time.Time `gorm:"column:read_date"`
	UnreadCount int        `gorm:"column:unread_count"`
}
