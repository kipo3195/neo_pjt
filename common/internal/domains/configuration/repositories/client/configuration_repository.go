package client

import (
	"common/internal/domains/configuration/models"
	"errors"

	"gorm.io/gorm"
)

type configurationRepository struct {
	db *gorm.DB
}

type ConfigurationRepository interface {
	GetConfigHash() (string, error)
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.WorksInfo{})
}

func NewConfigurationRepository(db *gorm.DB) ConfigurationRepository {
	return &configurationRepository{
		db: db,
	}
}

func (r *configurationRepository) GetConfigHash() (string, error) {

	var configHash string
	result := r.db.Model(&models.WorksInfo{}).
		Where("kind = ?", "config").
		Select("value").
		First(&configHash)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // 조회 결과 없음
		}
		return "", result.Error // 다른 DB 오류
	}

	return configHash, nil
}
