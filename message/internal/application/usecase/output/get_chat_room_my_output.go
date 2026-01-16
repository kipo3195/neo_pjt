package output

import "time"

type GetChatRoomMyOutput struct {
	RoomKey     string
	RoomType    string
	ReadDate    *time.Time
	UnreadCount int
}
