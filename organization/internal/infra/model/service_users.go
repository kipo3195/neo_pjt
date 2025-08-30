package model

type ServiceUsers struct {
	UserHash string `gorm:"column:user_hash;primaryKey;comment:pk"`
	UserId   string `gorm:"column:user_id;comment:사용자 ID"`
	UseYn    string `gorm:"column:use_yn;default:Y;comment:사용여부"`
}

func (ServiceUsers) TableName() string {
	return "service_users"
}

// 서비스 등록 사용자 (공통)
