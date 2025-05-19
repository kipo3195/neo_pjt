package models

// 테넌트 (works)의 연결 정보 상세

type ConnectInfo struct {
	Domain     string `json:"domain"`
	ApiVersion string `json:"api_version"`
	UdtDate    string `json:"udt_date"`
}

func (ConnectInfo) TableName() string {
	return "connect_info"
}
