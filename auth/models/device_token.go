package models

type DeviceToken struct {
	Seq   int    `gorm:"column:seq;primaryKey;autoIncrement"`
	Uuid  string `gorm:"column:uuid"`
	Token string `gorm:"column:token"`
}

// column:seq → 이 필드가 DB에서 seq라는 이름의 컬럼과 매핑됨을 명시적으로 지정.
// primaryKey → 기본 키임을 명시.
// autoIncrement → 자동 증가하는 필드임을 명시.
// json:"seq" → JSON으로 직렬화될 때 필드 이름을 seq로 지정.

func (DeviceToken) TableName() string {
	return "device_token"
}
