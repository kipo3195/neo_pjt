package model

// 인증 처리를 위한 정보
type WorksList struct {
	ServerUrl string `gorm:"column:server_url"` // 실제로 접속하는 url domain
	Code      string `gorm:"column:works_code"`
	Name      string `gorm:"column:works_name"`
	UseYn     string `gorm:"column:use_yn"`
	RegDate   string `gorm:"column:reg_date"`
}

func (WorksList) TableName() string {
	return "works_list"
}
