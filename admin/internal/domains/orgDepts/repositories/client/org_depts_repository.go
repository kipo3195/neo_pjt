package client

import "gorm.io/gorm"

type orgDeptsRepository struct {
	db *gorm.DB
}

type OrgDeptsRepository interface {
}

func NewOrgDeptsRepository(db *gorm.DB) OrgDeptsRepository {

	return &orgDeptsRepository{
		db: db,
	}

}
