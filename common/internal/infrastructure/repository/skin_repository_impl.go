package repository

import (
	"common/internal/domain/skin/entity"
	"common/internal/domain/skin/repository"
	"common/internal/infrastructure/model"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type skinRepositoryImpl struct {
	db *gorm.DB
}

func SkinMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.AppSkinConfig{})
	db.AutoMigrate(&model.AppSkinFileInfo{})
}

func NewSkinRepository(db *gorm.DB) repository.SkinRepository {
	return &skinRepositoryImpl{
		db: db,
	}
}

func (r *skinRepositoryImpl) GetSkinHash() (string, error) {

	var skinHash string
	result := r.db.Model(&model.AppSkinConfig{}).
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

func (r *skinRepositoryImpl) PutSkinFileInfo(ctx context.Context, entity *entity.SkinFileInfoEntity) (bool, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 스킨 해시 insert
	if err := tx.Create(&model.AppSkinConfig{
		Value: entity.FileHash,
	}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 스킨 정보 저장
	// insert 처리
	if err := tx.Create(&model.AppSkinFileInfo{
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
