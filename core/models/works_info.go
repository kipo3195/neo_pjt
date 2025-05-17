package models

// 인증 처리를 위한 정보
type WorksInfo struct {
	Code    string `gorm:"column:works_code"`
	Name    string `gorm:"column:works_name"`
	Domain  string `gorm:"column:works_domain"`
	UseYn   string `gorm:"column:use_yn"`
	RegDate string `gorm:"column:reg_date"`
}

func (WorksInfo) TableName() string {
	return "works_info"
}
