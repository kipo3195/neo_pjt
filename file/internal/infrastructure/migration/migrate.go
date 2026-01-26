package migration

import (
	"file/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.FileUrlMigrate(db)
}
