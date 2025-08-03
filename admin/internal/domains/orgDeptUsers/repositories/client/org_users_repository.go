package client

import "gorm.io/gorm"

type orgDeptUsersRepository struct {
	db *gorm.DB
}

type OrgDeptUsersRepository interface {
}

func NewOrgDeptUsersRepository(db *gorm.DB) OrgDeptUsersRepository {

	return &orgDeptUsersRepository{
		db: db,
	}

}
