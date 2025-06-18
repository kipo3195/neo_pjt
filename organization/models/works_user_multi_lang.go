package models

import "time"

type WorksUserMultiLang struct {
	UserHash   string    `gorm:"column:user_hash;primaryKey;comment:pk"`
	KrLang     string    `gorm:"column:kr_lang;comment:한국어"`
	EnLang     string    `gorm:"column:en_lang;comment:영어"`
	CnLang     string    `gorm:"column:cn_lang;comment:중국어"`
	JpLang     string    `gorm:"column:jp_lang;comment:일본어"`
	CreateDate time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

func (WorksUserMultiLang) TableName() string {
	return "works_user_multi_lang"
}
