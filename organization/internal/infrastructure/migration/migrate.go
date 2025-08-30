package migration

import (
	"org/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.UserMigrate(db)
	repository.OrgMigrate(db)
	repository.DepartmentMigrate(db)
}
