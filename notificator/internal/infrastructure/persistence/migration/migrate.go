package migration

import (
	"notificator/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.ChatMigrate(db)
	repository.ServiceUsersMigrate(db)
}
