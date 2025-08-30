package model

import "time"

type PositionMultiLang struct {
	PositionCode string    `gorm:"column:position_code;primaryKey;comment:직책 - 파트장, 팀장, 본부장"`
	KrLang       string    `gorm:"column:kr_lang;comment:한국어"`
	EnLang       string    `gorm:"column:en_lang;comment:영어"`
	CnLang       string    `gorm:"column:cn_lang;comment:중국어"`
	JpLang       string    `gorm:"column:jp_lang;comment:일본어"`
	CreateDate   time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

// 직위 (사원, 대리) 다국어 처리
func (PositionMultiLang) TableName() string {
	return "position_multi_lang"
}
