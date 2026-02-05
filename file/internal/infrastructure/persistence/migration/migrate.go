package migration

import (
	"file/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.FileUrlMigrate(db)
}
