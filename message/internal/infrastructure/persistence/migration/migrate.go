package migration

import (
	"message/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	repository.ServiceUsersMigrate(db)
	repository.OtpKeyMigrate(db)
	repository.ChatRoomMigrate(db)
	repository.ChatLineEventMigrate(db)
	repository.ChatRoomTitleMigrate(db)
	repository.ChatFileMigrate(db)
}
