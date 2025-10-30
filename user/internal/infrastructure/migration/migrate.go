package migration

import (
	"user/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.ProfileMigrate(db)
	repository.UserDetailMigrate(db)
}
