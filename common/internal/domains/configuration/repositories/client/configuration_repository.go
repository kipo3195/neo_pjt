package client

import "gorm.io/gorm"

type configurationRepository struct {
	db *gorm.DB
}

type ConfigurationRepository interface {
}

func NewConfigurationRepository(db *gorm.DB) ConfigurationRepository {
	return &configurationRepository{
		db: db,
	}
}
