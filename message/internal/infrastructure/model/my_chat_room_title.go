package model

import "time"

type ChatRoomTitle struct {
	Org        string    `gorm:"column:org;primaryKey;type:varchar(20);comment:사용자 org code"`
	UserHash   string    `gorm:"column:user_hash;primaryKey;type:varchar(191);comment:사용자 hash"`
	RoomKey    string    `gorm:"column:room_key;primaryKey;type:varchar(50);comment:고정된 형식의 룸 키"`
	MyTitle    string    `gorm:"column:my_room_title;type:varchar(100);comment:내가 설정한 방 제목"`
	UpdateFlag string    `gorm:"column:update_flag;type:varchar(3);not null;comment:변경 유무 Y인 경우 나의 title로, N인 경우 생성시 title로"`
	UpdateDate time.Time `gorm:"column:update_date;not null;comment:수정 시간"`
}
