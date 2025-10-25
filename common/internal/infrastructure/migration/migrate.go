package migration

import (
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.SkinMigrate(db)
	repository.WorksInfoMigrate(db)
	repository.ConfigurationMigrate(db)
	repository.ProfileMigrate(db)
}
