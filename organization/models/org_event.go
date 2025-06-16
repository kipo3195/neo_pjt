package models

type OrgEvent struct {
	Seq        int    `gorm:"column:seq;primaryKey;comment:'seq'"`
	EventType  string `gorm:"column:event_type;comment:'이벤트 타입 C, U, D'"`
	Kind       string `gorm:"column:kind;comment:'그룹 0, 사용자 1'"`
	OrgCode    string `gorm:"column:org_code;comment:'org code'"`
	UpdateHash string `gorm:"column:update_hash;comment:'update hash 정보'"`
}

func (OrgEvent) TableName() string {
	return "org_event"
}
