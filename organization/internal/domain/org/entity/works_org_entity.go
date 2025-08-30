package entity

type WorksOrg struct {
	Org            string `gorm:"column:org" json:"org"` // 추가
	DeptCode       string `gorm:"column:dept_code"        json:"dept_code"`
	ParentDeptCode string `gorm:"column:parent_dept_code" json:"parent_dept_code"`
	KoLang         string `gorm:"column:ko_lang"     json:"koLang"`
	EnLang         string `gorm:"column:en_lang"     json:"enLang"`
	ZhLang         string `gorm:"column:cn_lang"     json:"zhLang"`
	JpLang         string `gorm:"column:jp_lang"     json:"jpLang"`
	RuLang         string `gorm:"column:ru_lang"     json:"ruLang"`
	ViLang         string `gorm:"column:vi_lang"     json:"viLang"`
	UpdateHash     string `gorm:"column:update_hash" json:"update_hash"`
	Kind           string `gorm:"column:kind" json:"kind"` // 추가
	Id             string `gorm:"column:id" json:"id"`
	Header         string `gorm:"column:header" json:"header"`
}
