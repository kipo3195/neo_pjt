package models

// 테넌트 (works)의 연결 정보 상세

type ConnectInfo struct {
	Domain     string `gorm:"column:domain"`
	ApiVersion string `gorm:"column:api_version"`
	UdtDate    string `gorm:"column:udt_date"`
}

func (ConnectInfo) TableName() string {
	return "connect_info"
}
