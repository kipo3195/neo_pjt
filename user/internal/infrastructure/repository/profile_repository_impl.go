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
		UserHash:            entity.UserHash,
		UserId:              entity.UserId,
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

func (r *profileRepositoryImpl) DeleteUserProfileImgInfo(ctx context.Context, userHash string, fileName string) error {

	// 단일 UPDATE (트랜잭션 불필요)
	result := r.db.WithContext(ctx).
		Model(&model.ProfileImgInfo{}).
		Where("user_hash = ? AND save_name = ?", userHash, fileName).
		Update("use_yn", "N")

	if result.Error != nil {
		log.Printf("[DeleteUserProfileImgInfo] - Update failed: %v\n", result.Error)
		return consts.ErrProfileImgDBDeleteError
	}

	if result.RowsAffected == 0 {
		log.Printf("[DeleteUserProfileImgInfo] - No rows affected for userhash=%s, fileName=%s\n", userHash, fileName)
		return consts.ErrProfileImgDBDeleteError
	}

	log.Println("[DeleteUserProfileImgInfo] - Update success")
	return nil

}

func (r *profileRepositoryImpl) RollbackDeleteUserProfileImgInfo(ctx context.Context, userHash string, fileName string) error {
	// 단일 UPDATE (트랜잭션 불필요)
	result := r.db.WithContext(ctx).
		Model(&model.ProfileImgInfo{}).
		Where("user_hash = ? AND save_name = ?", userHash, fileName).
		Update("use_yn", "Y")

	if result.Error != nil {
		log.Printf("[RollbackDeleteUserProfileImgInfo] - Update failed: %v\n", result.Error)
		return consts.ErrProfileImgDBRoleBackError
	}

	if result.RowsAffected == 0 {
		log.Printf("[RollbackDeleteUserProfileImgInfo] - No rows affected for userhash=%s, fileName=%s\n", userHash, fileName)
		return consts.ErrProfileImgDBRoleBackError
	}

	log.Println("[RollbackDeleteUserProfileImgInfo] - Update success")
	return nil
}

func (r *profileRepositoryImpl) GetProfileInfo(ctx context.Context, en entity.GetProfileInfoEntity) (map[string]entity.GetProfileInfoResultEntity, error) {

	var model []model.ProfileImgInfo

	err := r.db.Table("profile_img_info AS p1").
		Select("p1.user_hash, p1.img_hash, p1.save_name, p1.save_path, p1.size, p1.create_at, p1.use_yn").
		Where("p1.user_hash IN ?", en.UserHash).
		Where("p1.create_at = (?)",
			r.db.Table("profile_img_info AS p2").
				Select("MAX(p2.create_at)").
				Where("p2.user_hash = p1.user_hash"),
		).
		Scan(&model).Error

	if err != nil {
		return nil, err
	}

	temp := make(map[string]entity.GetProfileInfoResultEntity)

	for i := 0; i < len(model); i++ {
		e := entity.GetProfileInfoResultEntity{
			UserHash:    model[i].UserHash,
			ProfileHash: model[i].ProfileImgHash,
		}
		temp[model[i].UserHash] = e
	}

	return temp, nil
}
