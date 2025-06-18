package models

type OrgEvent struct {
	Seq        int    `gorm:"column:seq;primaryKey;comment:seq"`
	EventType  string `gorm:"column:event_type;comment:이벤트 타입 C, U, D"`
	Id         string `gorm:"column:id; comment:부서면 부서코드, 사용자면 사용자 hash"`
	Kind       string `gorm:"column:kind;comment:부서 0, 사용자 1"`
	OrgCode    string `gorm:"column:org_code;comment:org code"`
	UpdateHash string `gorm:"column:update_hash;comment:update hash 정보"` // 필요 없을 수도
}

func (OrgEvent) TableName() string {
	return "org_event"
}
