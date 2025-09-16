package repository

import (
	"auth/internal/consts"
	"auth/internal/domain/userAuth/entity"
	"auth/internal/domain/userAuth/repository"
	"auth/internal/infrastructure/model"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type userAuthRepository struct {
	db *gorm.DB
}

func NewUserAuthRepository(db *gorm.DB) repository.UserAuthRepository {
	return &userAuthRepository{
		db: db,
	}
}

func UserAuthMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.UserAuth{})
}

func (r *userAuthRepository) PutUserAuthInfo(ctx context.Context, entity entity.UserAuthInfoEntity) error {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&model.UserAuth{
		Id:       entity.Id,
		Salt:     entity.Salt,
		UserHash: entity.UserHash,
		AuthHash: entity.AuthHash,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutUserAuth] - Commit failed")
		return err
	}
	log.Println("[PutUserAuth] - Commit Success")
	return nil
}

func (r *userAuthRepository) GetUserSalt(ctx context.Context, id string) (string, error) {
	var salt string
	result := r.db.Model(&model.UserAuth{}).
		Where("id = ?", id).
		Select("salt").
		First(&salt)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", consts.ErrSaltNotRegist // 조회 결과 없음
		}
		return "", result.Error // 다른 DB 오류
	}

	return salt, nil
}

func (r *userAuthRepository) GetUserAuthHash(ctx context.Context, id string) (string, error) {
	var authHash string
	result := r.db.Model(&model.UserAuth{}).
		Where("id = ?", id).
		Select("auth_hash").
		First(&authHash)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", result.Error // 조회 결과 없음
		}
		return "", result.Error // 다른 DB 오류
	}
	return authHash, nil
}
