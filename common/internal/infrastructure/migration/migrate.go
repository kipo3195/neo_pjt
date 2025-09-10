package migration

import (
	configuration "common/internal/domains/configuration/repositories/client"
	"common/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.SkinMigrate(db)
	repository.WorksInfoMigrate(db)
	configuration.Migrate(db)
}
