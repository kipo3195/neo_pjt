package model

type ChatRoom struct {
	RoomKey  string `gorm:"column:room_key;type:varchar(50);primaryKey;comment:고정된 생성 형식. 동일한 참여자 생성의 경우 체크할 수 있도록 함"`
	RoomType string `gorm:"column:room_type;type:varchar(10);comment:일반 : N 오픈 : O"`
}

func (ChatRoom) TableName() string {
	return "chat_room"
}
