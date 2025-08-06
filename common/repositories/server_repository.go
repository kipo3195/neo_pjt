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
	GetConnectInfo(worksCode string) (*entities.ConnectInfo, error)

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
