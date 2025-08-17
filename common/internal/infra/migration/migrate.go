package migration

import (
	configuration "common/internal/domains/configuration/repositories/client"
	device "common/internal/domains/device/repositories/serverRepository"
	skin "common/internal/domains/skin/repositories/server"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	skin.Migrate(db)
	device.Migrate(db)
	configuration.Migrate(db)
}
