package server

import (
	"common/models"
	"errors"

	"gorm.io/gorm"
)

type skinRepository struct {
	db *gorm.DB
}

type SkinRepository interface {
	GetSkinHash() (string, error)
}

func NewSkinRepository(db *gorm.DB) SkinRepository {
	return &skinRepository{
		db: db,
	}
}

func (r *skinRepository) GetSkinHash() (string, error) {

	var skinHash string
	result := r.db.Model(&models.AppSkinConfig{}).
		Where("kind = ?", "skin").
		Select("value").
		First(&skinHash)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // 조회 결과 없음
		}
		return "", result.Error // 다른 DB 오류
	}

	return skinHash, nil
}
