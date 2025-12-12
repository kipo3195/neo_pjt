package repository

import (
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/model"

	"gorm.io/gorm"
)

type orgInfoBatchRepository struct {
}

func OrgInfoMigrate(db *gorm.DB) {
	db.AutoMigrate(model.OrgInfo{})
}

func NewOrgInfoBatchRepository() repository.OrgInfoBatchRepository {

	return &orgInfoBatchRepository{}

}
