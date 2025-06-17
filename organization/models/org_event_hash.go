package models

type OrgEventHash struct {
	Seq        int    `gorm:"column:seq;primaryKey;comment:'seq'"`
	OrgCode    string `gorm:"column:org_code;comment:'org code'"`
	UpdateHash string `gorm:"column:update_hash;comment:'update hash 정보'"`
	adminHash  string `gorm:"column:admin_hash;comment'작업한 관리자의 hash 정보'"`
}

func (OrgEventHash) TableName() string {
	return "org_event_hash"
}
