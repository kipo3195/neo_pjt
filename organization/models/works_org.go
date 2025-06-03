package models

// 부서 - 사용자 재귀 쿼리 조회 매핑 정보
type WorksOrg struct {
	DeptCode       string `gorm:"column:dept_code"        json:"dept_code"`
	ParentDeptCode string `gorm:"column:parent_dept_code" json:"parent_dept_code"`
	KrLang         string `gorm:"column:kr_lang"     json:"krLang"`
	EnLang         string `gorm:"column:en_lang"     json:"enLang"`
	CnLang         string `gorm:"column:cn_lang"     json:"cnLang"`
	JpLang         string `gorm:"column:jp_lang"     json:"jpLang"`
	UpdateHash     string `gorm:"column:update_hash" json:"update_hash"`
}
