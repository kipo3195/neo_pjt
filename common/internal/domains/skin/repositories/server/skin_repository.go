package server

import (
	"common/internal/domains/skin/entities"
	"common/internal/domains/skin/models"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type skinRepository struct {
	db *gorm.DB
}

type SkinRepository interface {
	GetSkinHash() (string, error)
	PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error)
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.AppSkinConfig{})
	db.AutoMigrate(&models.AppSkinFileInfo{})
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

func (r *skinRepository) PutSkinFileInfo(ctx context.Context, entity *entities.SkinFileInfoEntity) (bool, error) {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 스킨 해시 insert
	if err := tx.Create(&models.AppSkinConfig{
		Value: entity.FileHash,
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
