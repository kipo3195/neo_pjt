package migration

import (
	"admin/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.ServiceUsersMigrate(db)
}
