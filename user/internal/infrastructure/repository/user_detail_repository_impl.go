package repository

import (
	"context"
	"user/internal/domain/userDetail/entity"
	"user/internal/domain/userDetail/repository"
	"user/internal/infrastructure/model"

	"gorm.io/gorm"
)

type userDetailRepositoryImpl struct {
	db *gorm.DB
}

func NewUserDetailRepository(db *gorm.DB) repository.UserDetailRepository {
	return &userDetailRepositoryImpl{
		db: db,
	}
}

func UserDetailMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.UserDetail{})
}

func (r *userDetailRepositoryImpl) GetUserInfoDetailInfo(ctx context.Context, entity entity.GetUserDetailInfoEntity) ([]entity.UserDetailInfoEntity, error) {

	return nil, nil
}
