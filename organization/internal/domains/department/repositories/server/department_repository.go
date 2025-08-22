package server

import "gorm.io/gorm"

type departmentRepository struct {
	db *gorm.DB
}

type DepartmentRepository interface {
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}
