package chatRoom

import "time"

type GetChatRoomMyDto struct {
	RoomKey     string     `json:"roomKey"`
	RoomType    string     `json:"roomType"`
	UnreadCount int        `json:"unreadCount"`
	ReadDate    *time.Time `json:"readDate,omitempty"`
}
