package serverRepository

import (
	"common/internal/domains/device/entities"
	"common/internal/domains/device/models"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

type DeviceRepository interface {
	GetConnectInfo(worksCode string) (*entities.ConnectInfo, error)
	GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error)
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{
		db: db,
	}
}

func (r *deviceRepository) GetConnectInfo(worksCode string) (*entities.ConnectInfo, error) {

	// model
	var connectInfo models.ConnectInfo

	// domain으로 auth에 접근할 것인가?
	result := r.db.Where("works_code = ?", worksCode).First(&connectInfo)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetConnectInfo] - DB error")
		return nil, result.Error
	} else {

		if result.RowsAffected > 0 {
			return &entities.ConnectInfo{
				ServerUrl: connectInfo.ServerUrl,
			}, nil
		} else {
			log.Println("[GetConnectInfo] - DB select X")
			return nil, nil
		}
	}

}

func (r *deviceRepository) GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return &entities.WorksConfig{}, tx.Error
	}

	// 스킨 정보 조회
	var appSkinInfo models.AppSkinConfig
	err := tx.First(&appSkinInfo).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 진짜 DB 오류 → 롤백 후 에러 반환
		tx.Rollback()
		return nil, err
	}

	// 스킨 파일의 정보 조회
	var appSkinFileInfo []models.AppSkinFileInfo
	if err := tx.Find(&appSkinFileInfo).Error; err != nil {
		// 진짜 DB 오류일 경우만 처리
		log.Println("조회 실패:", err)
		return nil, err
	}

	// 언어, 타임존, 서버 설정, 컬러 조회
	log.Println("entity.WorksCode ", entity.WorksCode)
	var worksInfo []models.WorksInfo
	if err := tx.Where("works_code = ?", entity.WorksCode).Find(&worksInfo).Error; err != nil {
		tx.Rollback()
		return nil, err // ← 실제 DB 에러일 때만 종료
	}

	log.Println("worksInfo : ", worksInfo)

	if err := tx.Commit().Error; err != nil {
		log.Println("[GetWorksConfig] - Commit failed")
		return &entities.WorksConfig{}, err
	}

	return toWorksWorksConfigEntity(appSkinInfo, worksInfo, appSkinFileInfo), nil
}

func toWorksWorksConfigEntity(appSkinConfig models.AppSkinConfig, worksInfo []models.WorksInfo, appSkinFileInfo []models.AppSkinFileInfo) *entities.WorksConfig {

	var timeZone, language, configHash string

	for _, info := range worksInfo {
		switch info.Kind {
		case "timeZone":
			timeZone = info.Value
		case "language":
			language = info.Value
		case "configHash":
			configHash = info.Value
		}
	}

	log.Println("toWorksConfigEntity")
	skinFileInfo := ConvertToSkinFileInfoEntities(appSkinFileInfo)
	return &entities.WorksConfig{
		TimeZone:   timeZone,
		Language:   language,
		ConfigHash: configHash,
		SkinHash:   appSkinConfig.Value,
		Skin:       skinFileInfo,
	}
}

func ConvertToSkinFileInfoEntities(models []models.AppSkinFileInfo) []entities.SkinFileInfoEntity {
	result := make([]entities.SkinFileInfoEntity, len(models))
	for i, m := range models {
		result[i] = entities.SkinFileInfoEntity{
			FileHash: m.FileHash,
			SkinType: m.SkinType,
		}
	}
	return result
}
