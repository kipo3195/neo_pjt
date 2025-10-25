package repository

import (
	"common/internal/consts"
	"common/internal/domain/profile/entity"
	"common/internal/domain/profile/repository"
	"common/internal/infrastructure/model"
	"context"
	"log"

	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) repository.ProfileRepository {

	return &profileRepositoryImpl{
		db: db,
	}
}

func ProfileMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ProfileImgInfo{})
}

func (r *profileRepositoryImpl) PutUserProfileImgInfo(ctx context.Context, entity entity.ProfileImgEntity) error {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&model.ProfileImgInfo{
		Id:                  entity.UserId,
		ProfileImgHash:      entity.ProfileImgHash,
		ProfileImgSavedName: entity.ProfileImgSavedName,
		ProfileImgSavedPath: entity.ProfileImgSavedPath,
		ProfileImgSize:      entity.ProfileImgSize,
	}).Error; err != nil {
		tx.Rollback()
		return consts.ErrProfileImgDBSaveError
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutUserProfileImgInfo] - Commit failed")
		return err
	}
	log.Println("[PutUserProfileImgInfo] - Commit Success")

	return nil
}
