package models

// 인증 처리를 위한 정보

type AuthInfo struct {
	Id       string `gorm:"column:id"`
	Password string `gorm:"column:password"`
}

func (AuthInfo) TableName() string {
	return "auth_info"
}
