package migration

import (
	"core/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.AppValidationMigrate(db)
}
