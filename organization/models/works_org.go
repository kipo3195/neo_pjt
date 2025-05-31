package models

type WorksOrg struct {
	DeptCode       string `gorm:"column:dept_code"        json:"dept_code"`
	ParentDeptCode string `gorm:"column:parent_dept_code" json:"parent_dept_code"`
	KrLang         string `gorm:"column:kr_lang"     json:"krLang"`
	EnLang         string `gorm:"column:en_lang"     json:"enLang"`
	CnLang         string `gorm:"column:cn_lang"     json:"cnLang"`
	JpLang         string `gorm:"column:jp_lang"     json:"jpLang"`
	DeptUpdateHash string `gorm:"column:dept_update_hash" json:"dept_update_hash"`
}
