package migration

import (
	"batch/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.OrgInfoMigrate(db)
}
