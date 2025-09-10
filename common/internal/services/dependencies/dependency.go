package dependencies

import (
	"common/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type Dependency struct {
	DB                *gorm.DB
	ConfigHashStorage storage.ConfigHashStorage
	SkinStorage       storage.SkinStorage
	AutoMigrate       bool
}
