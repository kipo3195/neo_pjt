package model

import "time"

type ChatRoomMember struct {
	RoomKey         string    `gorm:"column:room_key;type:varchar(50);primaryKey;comment:고정된 생성 형식. 동일한 참여자 생성의 경우 체크할 수 있도록 함"`
	MemberHash      string    `gorm:"column:member_hash;primaryKey;not null;comment:참여자 hash"`
	MemberState     string    `gorm:"column:member_state;type:varchar(1);not null;default:1;comment:참여자 상태 1 : 참여중, 2 : 퇴장, 3 : 강퇴"`
	MemberFirstDate time.Time `gorm:"column:member_first_date;not null;comment:참여자 최초 입장 시간"`
	MemberDate      time.Time `gorm:"column:member_date;not null;comment:참여자 입장 시간"`
	MemberExitDate  time.Time `gorm:"column:member_exit_date;comment:참여자 퇴장(강퇴) 시간"`
	MemberWorksCode string    `gorm:"column:member_works_code;type:varchar(50);comment:방 참여자 기준 works code"`
}

func (ChatRoomMember) TableName() string {
	return "chat_room_member"
}

// 20260101 정리
