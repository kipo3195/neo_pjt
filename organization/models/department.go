package models

type Department struct {
	DeptName string `gorm:"column:dept_name"`
}

func (Department) TableName() string {
	return "department"
}
