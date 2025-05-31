package models

type WorksDeptMultiLang struct {
	Seq      int    `gorm:"column:seq;primaryKey;autoIncrement;comment:'pk'"`
	DeptCode string `gorm:"column:dept_code;comment:'부서 코드 - 다국어 매핑용'"`
	DeptOrg  string `gorm:"column:dept_org;comment:'부서 코드 - 다국어 매핑용'"`
	KrLang   string `gorm:"column:kr_lang;comment:'한국어'"`
	EnLang   string `gorm:"column:en_lang;comment:'영어'"`
	CnLang   string `gorm:"column:cn_lang;comment:'중국어'"`
	JpLang   string `gorm:"column:jp_lang;comment:'일본어'"`
}

func (WorksDeptMultiLang) TableName() string {
	return "works_dept_multi_lang"
}
