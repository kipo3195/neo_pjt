package model

type OrgEventHash struct {
	Seq        int    `gorm:"column:seq;primaryKey;comment:seq"`
	OrgCode    string `gorm:"column:org_code;comment:org code"`
	UpdateHash string `gorm:"column:update_hash;comment:update hash 정보"`
	AdminHash  string `gorm:"column:admin_hash;comment작업한 관리자의 hash 정보"`
}

func (OrgEventHash) TableName() string {
	return "org_event_hash"
}

// 각 org별 최신의 hash정보를 관리하는 테이블
