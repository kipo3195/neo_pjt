package repository

import (
	"common/internal/domain/org/repository"
	"common/internal/infrastructure/model"

	"gorm.io/gorm"
)

type orgRepositoryImpl struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) repository.OrgRepository {
	return &orgRepositoryImpl{
		db: db,
	}
}

func OrgMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.OrgCode{})
}

func (r *orgRepositoryImpl) GetOrgCode() ([]string, error) {
	var orgCodes []string

	err := r.db.
		Raw("SELECT org FROM org_code").
		Scan(&orgCodes).Error

	if err != nil {
		return nil, err

	}
	return orgCodes, nil
}
