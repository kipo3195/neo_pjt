package migration

import (
	"auth/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.TokenMigrate(db)
	repository.CertificationMigrate(db)
	repository.UserAuthMigrate(db)
}
