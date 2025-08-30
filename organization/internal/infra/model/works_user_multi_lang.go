package model

import "time"

type WorksUserMultiLang struct {
	UserHash   string    `gorm:"column:user_hash;primaryKey;comment:pk"`
	DefLang    string    `gorm:"column:def_lang;comment:기본 언어"`
	KoLang     string    `gorm:"column:ko_lang;comment:한국어"`
	EnLang     string    `gorm:"column:en_lang;comment:영어"`
	ZhLang     string    `gorm:"column:zh_lang;comment:중국어"`
	JpLang     string    `gorm:"column:jp_lang;comment:일본어"`
	RuLang     string    `gorm:"column:ru_lang;comment:러시아어"`
	ViLang     string    `gorm:"column:vi_lang;comment:베트남어"`
	CreateDate time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

func (WorksUserMultiLang) TableName() string {
	return "works_user_multi_lang"
}
