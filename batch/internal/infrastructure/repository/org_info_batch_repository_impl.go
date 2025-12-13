package repository

import (
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/model"

	"gorm.io/gorm"
)

type orgInfoRepository struct {
}

func OrgInfoMigrate(db *gorm.DB) {
	db.AutoMigrate(model.OrgInfo{})
}

func NewOrgInfoRepository() repository.OrgInfoRepository {

	return &orgInfoRepository{}
}
