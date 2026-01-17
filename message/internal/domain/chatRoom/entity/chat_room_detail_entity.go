package entity

type ChatRoomDetailEntity struct {
	RoomKey         string `gorm:"column:room_key"`
	Title           string `gorm:"column:room_title"`
	LineKey         string `gorm:"column:line_key"`
	EventType       string `grom:"column:event_type"`
	Cmd             int    `gorm:"column:cmd"`
	Contents        string `gorm:"column:contents"`
	SendDate        string `gorm:"column:send_date"`
	SecretFlag      string `gorm:"column:room_secret_flag"`
	Secret          string `gorm:"column:room_secret"`
	Description     string `gorm:"column:room_description"`
	State           string `gorm:"column:room_state"`
	WorksCode       string `gorm:"column:room_works_code"`
	CreateDate      string `gorm:"column:room_create_date"`
	CreateUser      string `gorm:"column:room_create_user"`
	Member          string `gorm:"column:member"`
	Owner           string `gorm:"column:owner"`
	Type            string `gorm:"column:room_type"`
	MyRoomTitle     string `gorm:"column:my_room_title"`
	TitleUpdateDate string `gorm:"column:title_update_date"`
	TitleUpdateFlag string `gorm:"column:title_update_flag"`
	UnreadCount     int    `gorm:"column:unread_count"`
	UnreadCountDate string `gorm:"column:unread_count_date"`
	LastReadDate    string `gorm:"column:last_read_date"`
}
