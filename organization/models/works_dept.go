package models

import "time"

type WorksDept struct {
	DeptCode        string    `gorm:"column:dept_code;primaryKey;size:20;comment:'pk'"`
	DeptOrg         string    `gorm:"column:dept_org;comment:'부서의 계열사 코드'"`
	ParentsDeptCode string    `gorm:"column:parent_dept_code;comment:'부모 코드'"`
	DeptCreateDate  time.Time `gorm:"column:dept_create_date;default:CURRENT_TIMESTAMP;comment:'등록일'"`
	UpdateHash      string    `gorm:"column:update_hash;comment:'해시 정보'"`
	UseYn           string    `gorm:"column:use_yn;comment:'사용 여부'"`
}

func (WorksDept) TableName() string {
	return "works_dept"
}
