package repository

import (
	"common/internal/domain/user/repository"
	"context"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) CheckUserRegist(ctx context.Context, id string) error {

	// service users 조회 처리
	return nil
}
