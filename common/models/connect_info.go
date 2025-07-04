package models

// 테넌트 (works)의 연결 정보 상세

type ConnectInfo struct {
	ServerUrl  string `gorm:"column:server_url"`
	WorksCode  string `gorm:"column:works_code;primaryKey"`
	ApiVersion string `gorm:"column:api_version"`
	UdtDate    string `gorm:"column:udt_date"`
}

func (ConnectInfo) TableName() string {
	return "connect_info"
}
