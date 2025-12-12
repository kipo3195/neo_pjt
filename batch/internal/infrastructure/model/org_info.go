package model

type OrgInfo struct {
	Seq            int    `gorm:"column:seq;int(11);primaryKey;autoIncrement;comment:pk"`
	Org            string `gorm:"column:org;varchar(30);not null;comment:works code"`
	DeptCode       string `gorm:"column:dept_code;varchar(50);not null;comment:부서 코드"`
	ParentDeptCode string `gorm:"column:parent_dept_code;varchar(50);not null;comment:부모 부서 코드"`
	UserId         string `gorm:"column:user_id;varchar(200);comment:사용자 id"`
	KoLang         string `gorm:"column:ko_lang;varchar(100);comment:부서 혹은 사용자 이름 (한국어)"`
	Kind           int    `gorm:"column:kind;int(11);not null; comment:부서 : 0, 사용자 : 1"`
}

func (OrgInfo) TableName() string {
	return "org_info"
}
