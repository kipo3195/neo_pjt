package entity

type ChatRoomDetailEntity struct {
	RoomKey     string `gorm:"column:room_key"`
	Title       string `gorm:"column:room_title"`
	SecretFlag  string `gorm:"column:room_secret_flag"`
	Secret      string `gorm:"column:room_secret"`
	Description string `gorm:"column:room_description"`
	State       string `gorm:"column:room_state"`
	WorksCode   string `gorm:"column:room_works_code"`
	CreateDate  string `gorm:"column:room_create_date"`
	CreateUser  string `gorm:"column:room_create_user"`
	Hash        string `gorm:"column:room_hash"`
	Member      string `gorm:"column:member"`
}
