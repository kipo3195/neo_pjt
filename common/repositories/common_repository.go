package repositories

import (
	"common/entities"
	"common/models"
	"errors"

	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

type CommonRepository interface {
	GetSkinHash() (string, error)
	GetConfigHash() (string, error)
}

func NewCommonRepository(db *gorm.DB) CommonRepository {

	return &commonRepository{db: db}
}

func (r *commonRepository) GetSkinHash() (string, error) {

	var config models.AppSkinConfig
	result := r.db.Model(&models.AppSkinConfig{}).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // 조회 결과 없음은 에러 아님
		}
		return "", result.Error // 진짜 DB 오류
	}

	return config.SkinHash, nil
}

func (r *commonRepository) GetConfigHash() (string, error) {

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

func (r *commonRepository) GetSkinInfo() ([]entities.SkinFileInfoEntity, error) {

	var appSkinFileInfo []models.AppSkinFileInfo

	if err := r.db.Find(&appSkinFileInfo).Error; err != nil {
		return nil, err
	}

	entityList := toSkinFileInfoEntityList(appSkinFileInfo)
	return entityList, nil
}

func toSkinFileInfoEntityList(models []models.AppSkinFileInfo) []entities.SkinFileInfoEntity {
	result := make([]entities.SkinFileInfoEntity, 0, len(models))
	for _, m := range models {
		result = append(result, entities.SkinFileInfoEntity{
			FileHash: m.FileHash,
			SkinType: m.SkinType,
			FilePath: m.FilePath,
		})
	}
	return result
}
