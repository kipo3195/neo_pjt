package model

// 인증 처리를 위한 정보

type AuthInfo struct {
	Id       string `gorm:"column:id"`
	Password string `gorm:"column:password"`
	Userhash string `gorm:"column:user_hash"` // service_users의 join값.
}

func (AuthInfo) TableName() string {
	return "auth_info"
}
