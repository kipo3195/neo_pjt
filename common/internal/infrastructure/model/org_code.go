package model

type OrgCode struct {
	WorksCode  string `gorm:"column:works_code;type:varchar(30);primaryKey;comment:works code"`
	Org        string `gorm:"column:org;type:varchar(30);primaryKey;comment:org code"`
	CreateDate string `gorm:"column:create_date;comment:추가시간"`
}

func (OrgCode) TableName() string {
	return "org_code"
}
