package migration

import (
	configuration "common/internal/domains/configuration/repositories/client"
	skin "common/internal/domains/skin/repositories/server"
	worksInfo "common/internal/domains/worksInfo/repositories/serverRepository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	skin.Migrate(db)
	worksInfo.Migrate(db)
	configuration.Migrate(db)
}
