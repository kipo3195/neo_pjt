package models

import "time"

type WorksDept struct {
	Seq             int       `gorm:"column:seq;primaryKey;autoIncrement;comment:'pk'"`
	DeptOrg         string    `gorm:"column:dept_org;comment:'부서의 계열사 코드'"`
	DeptCode        string    `gorm:"column:dept_code;comment:'부서 코드 - 다국어 매핑용'"`
	ParentsDeptCode string    `gorm:"column:parent_dept_code;comment:'부모 코드'"`
	DeptCreateDate  time.Time `gorm:"column:dept_create_date;comment:'등록일'"`
	DeptUpdateHash  time.Time `gorm:"column:dept_update_hash;comment:'해시 정보'"`
	UseYn           string    `gorm:"column:use_yn;comment:'사용 여부'"`
}

func (WorksDept) TableName() string {
	return "works_dept"
}
