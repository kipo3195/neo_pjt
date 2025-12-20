package model

import "time"

type ChatRoomOwner struct {
	RoomKey         string    `gorm:"column:room_key;type:varchar(50);primaryKey;comment:고정된 생성 형식. 동일한 참여자 생성의 경우 체크할 수 있도록 함"`
	OwnerHash       string    `gorm:"column:owner_hash;type:varchar(191);primaryKey;not null;comment:방장의 user hash"`
	UpdateDate      time.Time `gorm:"column:update_date;not null;comment:방장 수정 시간"`
	ActiveFlag      string    `gorm:"column:active_flag;type:varchar(3);not null;comment:방장 활성화 여부"`
	MemberWorksCode string    `gorm:"column:member_works_code;type:varchar(50);comment:방 참여자 기준 works code"`
}

func (ChatRoomOwner) TableName() string {
	return "chat_room_owner"
}
