package repository

import (
	"common/internal/domain/configuration/repository"
	"common/internal/infrastructure/model"
	"errors"

	"gorm.io/gorm"
)

type configurationRepositoryImpl struct {
	db *gorm.DB
}

func ConfigurationMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.WorksInfo{})
}

func NewConfigurationRepository(db *gorm.DB) repository.ConfigurationRepository {
	return &configurationRepositoryImpl{
		db: db,
	}
}

func (r *configurationRepositoryImpl) GetConfigHash() (string, error) {

	var configHash string
	result := r.db.Model(&model.WorksInfo{}).
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
