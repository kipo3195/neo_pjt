package repository

import (
	"gorm.io/gorm"
)

type configurationRepository struct {
	db *gorm.DB
}

type ConfigurationRepository interface {
	GetConfigHash() (string, error)
}
