package migration

import (
	"notificator/internal/infrastructure/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.ChatMigrate(db)
}
