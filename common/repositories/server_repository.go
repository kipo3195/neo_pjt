package repositories

import (
	"common/entities"
	"common/models"
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

type ServerRepository interface {
	GetConnectInfo(worksCode string) (*entities.ConnectInfo, error)
	GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error)
	PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error)
}

func NewServerRepository(db *gorm.DB) ServerRepository {

	return &serverRepository{db: db}
}

func (r *serverRepository) GetConnectInfo(worksCode string) (*entities.ConnectInfo, error) {

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

func (r *serverRepository) GetWorksConfig(entity entities.GetWorksConfig, ctx context.Context) (*entities.WorksConfig, error) {

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
	fmt.Println("entity.WorksCode ", entity.WorksCode)
	var worksInfo []models.WorksInfo
	if err := tx.Where("works_code = ?", entity.WorksCode).Find(&worksInfo).Error; err != nil {
		tx.Rollback()
		return nil, err // ← 실제 DB 에러일 때만 종료
	}

	fmt.Println("worksInfo : ", worksInfo)

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

	fmt.Println("toWorksConfigEntity")
	skinFileInfo := ConvertToSkinFileInfoEntities(appSkinFileInfo)
	return &entities.WorksConfig{
		TimeZone:   timeZone,
		Language:   language,
		ConfigHash: configHash,
		SkinHash:   appSkinConfig.SkinHash,
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

func (r *serverRepository) PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 스킨 해시 insert
	if err := tx.Create(&models.AppSkinConfig{
		SkinHash: entity.FileHash,
	}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 스킨 정보 저장
	// insert 처리
	if err := tx.Create(&models.AppSkinFileInfo{
		SkinType: entity.SkinType,
		//FileUrl:  entity.FileUrl,
		//FileName: entity.FileName,
		//FilePath: entity.FilePath,
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
