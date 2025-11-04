package model

import "time"

type RoleMultiLang struct {
	RoleCode   string    `gorm:"column:role_code;primaryKey;comment:직책 - 파트장, 팀장, 본부장"`
	KrLang     string    `gorm:"column:kr_lang;comment:한국어"`
	EnLang     string    `gorm:"column:en_lang;comment:영어"`
	ZhLang     string    `gorm:"column:zh_lang;comment:중국어"`
	JpLang     string    `gorm:"column:jp_lang;comment:일본어"`
	CreateDate time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

// 직책 (파트장, 팀장) 다국어 처리
func (RoleMultiLang) TableName() string {
	return "role_multi_lang"
}
