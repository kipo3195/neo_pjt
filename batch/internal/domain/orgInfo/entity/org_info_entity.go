package entity

type OrgInfoEntity struct {
	Org            string `gorm:"column:org" json:"org"` // 추가
	DeptCode       string `gorm:"column:dept_code"        json:"deptCode"`
	ParentDeptCode string `gorm:"column:parent_dept_code" json:"parentDeptCode"`
	KoLang         string `gorm:"column:ko_lang"     json:"koLang"`
	EnLang         string `gorm:"column:en_lang"     json:"enLang"`
	ZhLang         string `gorm:"column:cn_lang"     json:"zhLang"`
	JpLang         string `gorm:"column:jp_lang"     json:"jpLang"`
	RuLang         string `gorm:"column:ru_lang"     json:"ruLang"`
	ViLang         string `gorm:"column:vi_lang"     json:"viLang"`
	UpdateHash     string `gorm:"column:update_hash" json:"updateHash"`
	Kind           string `gorm:"column:kind" json:"kind"` // 추가
	UserHash       string `gorm:"column:user_hash" json:"userHash"`
	UserId         string `gorm:"column:user_id" json:"userId"`
	Header         string `gorm:"column:header" json:"header"`
	Description    string `gorm:"column:description" json:"description"`
}
