package models

type ServiceUsers struct {
	UserHash string `gorm:"column:user_hash;primaryKey;comment:'pk'"`
	UserId   string `gorm:"column:user_id;commont:'사용자 ID'"`
}

func (ServiceUsers) TableName() string {
	return "service_users"
}

// 서비스 등록 사용자 (공통)
