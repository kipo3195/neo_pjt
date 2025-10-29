package repository

import (
	"context"
	"log"
	"user/internal/consts"
	"user/internal/domain/profile/entity"
	"user/internal/domain/profile/repository"
	"user/internal/infrastructure/model"

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

	err := r.db.WithContext(ctx).Create(&model.ProfileImgInfo{
		Id:                  entity.UserId,
		ProfileImgHash:      entity.ProfileImgHash,
		ProfileImgSavedName: entity.ProfileImgSavedName,
		ProfileImgSavedPath: entity.ProfileImgSavedPath,
		ProfileImgSize:      entity.ProfileImgSize,
	}).Error

	if err != nil {
		log.Printf("[PutUserProfileImgInfo] - DB insert failed: %v\n", err)
		return consts.ErrProfileImgDBSaveError
	}

	log.Println("[PutUserProfileImgInfo] - Insert Success")
	return nil
}

func (r *profileRepositoryImpl) DeleteUserProfileImgInfo(ctx context.Context, userId string, fileName string) error {

	// 단일 UPDATE (트랜잭션 불필요)
	result := r.db.WithContext(ctx).
		Model(&model.ProfileImgInfo{}).
		Where("id = ? AND save_name = ?", userId, fileName).
		Update("use_yn", "N")

	if result.Error != nil {
		log.Printf("[DeleteUserProfileImgInfo] - Update failed: %v\n", result.Error)
		return consts.ErrProfileImgDBDeleteError
	}

	if result.RowsAffected == 0 {
		log.Printf("[DeleteUserProfileImgInfo] - No rows affected for id=%s, fileName=%s\n", userId, fileName)
		return consts.ErrProfileImgDBDeleteError
	}

	log.Println("[DeleteUserProfileImgInfo] - Update success")
	return nil

}

func (r *profileRepositoryImpl) RollbackDeleteUserProfileImgInfo(ctx context.Context, userId string, fileName string) error {
	// 단일 UPDATE (트랜잭션 불필요)
	result := r.db.WithContext(ctx).
		Model(&model.ProfileImgInfo{}).
		Where("id = ? AND save_name = ?", userId, fileName).
		Update("use_yn", "Y")

	if result.Error != nil {
		log.Printf("[RollbackDeleteUserProfileImgInfo] - Update failed: %v\n", result.Error)
		return consts.ErrProfileImgDBRoleBackError
	}

	if result.RowsAffected == 0 {
		log.Printf("[RollbackDeleteUserProfileImgInfo] - No rows affected for id=%s, fileName=%s\n", userId, fileName)
		return consts.ErrProfileImgDBRoleBackError
	}

	log.Println("[RollbackDeleteUserProfileImgInfo] - Update success")
	return nil
}
