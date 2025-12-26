package model

type ChatRoomMember struct {
	RoomKey     string `gorm:"column:room_key;type:varchar(50);primaryKey;comment:고정된 생성 형식. 동일한 참여자 생성의 경우 체크할 수 있도록 함"`
	MemberHash  string `gorm:"column:member_hash;varchar(191);primaryKey;not null;comment:참여자 hash"`
	MemberState string `gorm:"column:member_state;type:varchar(1);not null;default:1;comment:참여자 상태 1 : 참여중, 2 : 퇴장, 3 : 강퇴"`
}

func (ChatRoomMember) TableName() string {
	return "chat_room_member"
}
