package model

import "time"

type WorksDept struct {
	DeptCode        string    `gorm:"column:dept_code;type:varchar(30);primaryKey;comment:pk"`
	DeptOrg         string    `gorm:"column:dept_org;type:varchar(10);primaryKey;comment:부서의 계열사 코드"`
	ParentsDeptCode string    `gorm:"column:parent_dept_code;type:varchar(30);comment:부모 부서 코드"`
	DeptCreateDate  time.Time `gorm:"column:dept_create_date;default:CURRENT_TIMESTAMP;comment:등록일"`
	UpdateHash      string    `gorm:"column:update_hash;type:varchar(30);comment:해시 정보"`
	Header          string    `gorm:"column:header;type:varchar(30);comment:부서 장"`
	Description     string    `gorm:"column:description;type:varchar(400);comment:부서 설명"`
	UseYn           string    `gorm:"column:use_yn;type:varchar(1);default:Y;comment:사용 여부"`
}

func (WorksDept) TableName() string {
	return "works_dept"
}
