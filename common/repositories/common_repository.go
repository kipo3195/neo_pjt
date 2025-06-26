package repositories

import (
	"common/entities"
	"common/models"

	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

type CommonRepository interface {
	GetSkinHashs() ([]entities.ConfigHashEntity, error)
	GetConfig() (entities.ConfigHashEntity, error)
}

func NewCommonRepository(db *gorm.DB) CommonRepository {

	return &commonRepository{db: db}
}

func (r *commonRepository) GetSkinHashs() ([]entities.ConfigHashEntity, error) {

	var appSkinConfig []models.AppSkinConfig
	result := r.db.Model(&models.AppSkinConfig{}).Scan(&appSkinConfig)

	if result.Error != nil {
		return nil, result.Error // 진짜 DB 오류
	}

	if result.RowsAffected == 0 {
		return nil, nil // 조회 결과 없음은 에러 아님
	}

	return toSkinHashEntity(appSkinConfig), nil
}

func toSkinHashEntity(models []models.AppSkinConfig) []entities.ConfigHashEntity {
	result := make([]entities.ConfigHashEntity, 0, len(models))

	for _, m := range models {
		result = append(result, entities.ConfigHashEntity{
			Device:   m.Device,
			SkinHash: m.Version,
		})
	}

	return result
}

func (r *commonRepository) GetConfig() (entities.ConfigHashEntity, error) {

	var config models.WorksInfo
	result := r.db.Where("kind = ?", "config").Model(&models.WorksInfo{}).Scan(&config)

	if result.Error != nil {
		return entities.ConfigHashEntity{}, result.Error // 진짜 DB 오류
	}

	if result.RowsAffected == 0 {
		return entities.ConfigHashEntity{}, nil // 조회 결과 없음은 에러 아님
	}

	return toConfigHashEntity(config), nil
}

func toConfigHashEntity(m models.WorksInfo) entities.ConfigHashEntity {
	return entities.ConfigHashEntity{
		Device:     m.Kind,
		ConfigHash: m.Value,
	}
}
