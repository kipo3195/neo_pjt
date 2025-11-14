package model

import "time"

type ProfileMsgInfo struct {
	UserHash       string    `gorm:"column:user_hash;type:varchar(100);primaryKey;	comment:사용자 hash 정보"`
	UserId         string    `gorm:"column:user_id;type:varchar(100);comment:사용자 계정"`
	ProfileMsg     string    `gorm:"column:profile_msg;type:mediumtext;not null;comment:프로필 메시지 정보"`
	CreateAt       time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
	ProfileMsgHash string    `gorm:"column:msg_hash;type:varchar(100);primaryKey;comment:프로필 hash 정보"`
}

func (ProfileMsgInfo) TableName() string {
	return "profile_msg_info"
}

// profileMsgHash가 필요한 이유?
// 만약 추후 멀티 프로필을 지원할때
// 특정 사용자 보여줄 프로필 매핑 테이블에서 활용
// 매핑 테이블은 등록 사용자, 타겟 사용자 profileImgInfo의 profileImgHash와 ProfileMsgInfo ProfileMsgHash를 가지며
// 등록사용자의hash를 key로 join하여 ProfileImgInfo와 ProfileMsgInfo의 기준으로 보여주도록 처리할 수 있다
