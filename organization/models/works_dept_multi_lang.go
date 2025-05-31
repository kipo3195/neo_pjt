package models

type WorksDeptMultiLang struct {
	Seq        int    `gorm:"column:seq;primaryKey;autoIncrement;comment:'pk'"`
	DeptCode   string `gorm:"column:dept_code;comment:'부서 코드 - 다국어 매핑용'"`
	DeptOrg    string `gorm:"column:dept_org;comment:'부서 코드 - 다국어 매핑용'"`
	DeptNameKr string `gorm:"column:dept_name_kr;comment:'한국어'"`
	DeptNameEn string `gorm:"column:dept_name_en;comment:'영어'"`
	DeptNameCn string `gorm:"column:dept_name_cn;comment:'중국어'"`
	DeptNameJp string `gorm:"column:dept_name_jp;comment:'일본어'"`
}

func (WorksDeptMultiLang) TableName() string {
	return "works_dept_multi_lang"
}
