package repository

import (
	"admin/internal/domain/orgDept/repository"

	"gorm.io/gorm"
)

type orgDeptRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgDeptRepository(db *gorm.DB) repository.OrgDeptRepository {

	return &orgDeptRepositoryImpl{
		db: db,
	}

}
