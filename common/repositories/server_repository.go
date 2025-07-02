package repositories

import (
	"common/entities"
	"common/models"
	"context"
	"log"

	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

type ServerRepository interface {
	GetConnectInfo(worksCode string) (*entities.InitResult, error)
	GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error)
	PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error)
}

func NewServerRepository(db *gorm.DB) ServerRepository {

	return &serverRepository{db: db}
}

func (r *serverRepository) GetConnectInfo(worksCode string) (*entities.InitResult, error) {

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

func (r *serverRepository) GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error) {

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

func (r *serverRepository) PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 스킨 해시 update
	if err := tx.Model(&models.AppSkinConfig{}).
		Where("1 = 1").
		Update("version", entity.FileHash).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 스킨 정보 저장
	// insert 처리
	if err := tx.Create(&models.AppSkinFileInfo{
		SkinType: entity.SkinType,
		Device:   entity.Device,
		FileName: entity.FileName,
		FilePath: entity.FilePath,
		FileHash: entity.FileHash,
	}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutSkinFileInfo] - Commit failed")
		return false, err
	}

	return false, nil
}
