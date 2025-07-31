package sharedModels

import "time"

type IssuedAppToken struct {
	Seq          int       `gorm:"column:seq;primaryKey;autoIncrement;comment:pk"`
	Uuid         string    `gorm:"column:uuid;type:varchar(100);comment:기기 고유값"`
	AppToken     string    `gorm:"column:app_token;type:varchar(400);comment:발급된 토큰 정보 JWT"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(400);comment:발급된 refresh 토큰 정보 JWT"`
	CreateAt     time.Time `gorm:"column:create_at;autoCreateTime;comment:DB 저장시간"`
}

// column:seq → 이 필드가 DB에서 seq라는 이름의 컬럼과 매핑됨을 명시적으로 지정.
// primaryKey → 기본 키임을 명시.
// autoIncrement → 자동 증가하는 필드임을 명시.
// json:"seq" → JSON으로 직렬화될 때 필드 이름을 seq로 지정.

func (IssuedAppToken) TableName() string {
	return "issued_app_tokens"
}
