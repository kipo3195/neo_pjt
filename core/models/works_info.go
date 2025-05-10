package models

// 인증 처리를 위한 정보

type WorksInfo struct {
	WCode   string `json:"w_code"`
	WName   string `json:"w_name"`
	WDomain string `json:"w_domain"`
	UseYn   string `json:"use_yn"`
	RegDate string `json:"reg_date"`
}

func (WorksInfo) TableName() string {
	return "works_info"
}
