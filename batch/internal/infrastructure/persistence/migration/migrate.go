package migration

import (
	"batch/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) error {
	err := repository.OrgInfoMigrate(db)
	if err != nil {
		return err
	}
	err = repository.UserDetailMigrate(db)
	if err != nil {
		return err
	}
	return nil
}
