package repository

import (
	"admin/internal/domain/orgDeptUser/repository"

	"gorm.io/gorm"
)

type orgDeptUsersRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgDeptUserRepository(db *gorm.DB) repository.OrgDeptUserRepository {

	return &orgDeptUsersRepositoryImpl{
		db: db,
	}

}
