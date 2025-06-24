package models

import "time"

type WorksDeptMultiLang struct {
	DeptCode   string    `gorm:"column:dept_code;primaryKey;comment:부서 코드 - 다국어 매핑용"`
	DeptOrg    string    `gorm:"column:dept_org;comment:부서 코드 - 다국어 매핑용"`
	KoLang     string    `gorm:"column:ko_lang;comment:한국어" json:"ko"`
	EnLang     string    `gorm:"column:en_lang;comment:영어" json:"en"`
	ZhLang     string    `gorm:"column:zh_lang;comment:중국어" json:"zh"`
	JpLang     string    `gorm:"column:jp_lang;comment:일본어" json:"jp"`
	RuLang     string    `gorm:"column:ru_lang;comment:러시아어" json:"ru"`
	ViLang     string    `gorm:"column:vi_lang;comment:베트남어" json:"vi"`
	DefLang    string    `gorm:"column:def_lang;comment:기본 언어" json:"def"`
	CreateDate time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
}

// 부서 다국어
func (WorksDeptMultiLang) TableName() string {
	return "works_dept_multi_lang"
}
