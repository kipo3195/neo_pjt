package repository

import (
	"auth/internal/domain/userAuth/entity"
	"auth/internal/domain/userAuth/repository"
	"auth/internal/infrastructure/model"
	"context"
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

func (r *userAuthRepository) PutUserAuth(ctx context.Context, entity entity.UserAuthEntity) error {

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 스킨 해시 insert
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
