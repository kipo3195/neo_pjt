package repositories

import (
	"common/entities"
	"common/models"
	"context"
	"log"

	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

type CommonRepository interface {
	GetConnectInfo(worksCode string) (*entities.InitResult, error)
	GetWorksConfig(GetWorksConfig entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error)
	GetSkinHashs() ([]entities.ConfigHashEntity, error)
	GetConfig() (entities.ConfigHashEntity, error)
}

func NewCommonRepository(db *gorm.DB) CommonRepository {

	return &commonRepository{db: db}
}

func (r *commonRepository) GetConnectInfo(worksCode string) (*entities.InitResult, error) {

	// model
	var connectInfo models.ConnectInfo

	// domain으로 auth에 접근할 것인가?
	result := r.db.Where("works_code = ?", worksCode).First(&connectInfo)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetConnectInfo] - DB error")
		return &entities.InitResult{}, result.Error
	} else {

		if result.RowsAffected > 0 {
			return &entities.InitResult{
				ConnectInfo: connectInfo.ServerUrl,
			}, nil
		} else {
			log.Println("[GetConnectInfo] - DB select X")
			return &entities.InitResult{}, nil
		}
	}

}

func (r *commonRepository) GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return &entities.WorksConfig{}, tx.Error
	}

	// 스킨 정보
	var appSkinInfo models.AppSkinConfig

	if err := tx.Where("works_code = ? AND device = ?", entity.WorksCode, entity.Device).First(&appSkinInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 언어, 타임존, 서버 설정 조회
	var worksInfo []models.WorksInfo

	if err := tx.Where("works_code = ?", entity.WorksCode).Find(&worksInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("[GetWorksConfig] - Commit failed")
		return &entities.WorksConfig{}, err
	}

	return toWorksConfigEntity(appSkinInfo, worksInfo), nil
}

func toWorksConfigEntity(appSkinConfig models.AppSkinConfig, worksInfo []models.WorksInfo) *entities.WorksConfig {

	var timeZone, language, configVersion string

	for _, info := range worksInfo {
		switch info.Kind {
		case "timeZone":
			timeZone = info.Value
		case "language":
			language = info.Value
		case "configVersion":
			configVersion = info.Value
		}
	}

	return &entities.WorksConfig{
		TimeZone:      timeZone,
		Language:      language,
		ConfigVersion: configVersion,
		SkinVersion:   appSkinConfig.Version,
	}
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
