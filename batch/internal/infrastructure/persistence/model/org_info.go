package model

type OrgInfo struct {
	Seq            int64  `gorm:"column:seq;primaryKey;autoIncrement;comment:pk" json:"seq"`
	Org            string `gorm:"column:org;type:longtext;not null;comment:org code" json:"org"`
	DeptCode       string `gorm:"column:dept_code;type:longtext;not null;comment:부서 코드" json:"dept_code"`
	ParentDeptCode string `gorm:"column:parent_dept_code;type:longtext;not null;comment:부모 부서 코드" json:"parent_dept_code"`
	UserId         string `gorm:"column:user_id;type:varchar(100);comment:사용자 id" json:"user_id,omitempty"`
	Kind           string `gorm:"column:kind;type:varchar(1);not null;comment:부서:0, 사용자:1" json:"kind"`
	KoLang         string `gorm:"column:ko_lang;type:varchar(800)" json:"ko_lang,omitempty"`
	UserHash       string `gorm:"column:user_hash;type:varchar(100);comment:사용자 id"`
}

func (OrgInfo) TableName() string {
	return "org_info"
}
