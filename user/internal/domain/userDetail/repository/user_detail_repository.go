package repository

import (
	"context"
	"user/internal/domain/userDetail/entity"

	"gorm.io/gorm"
)

type userDetailRepository struct {
	db *gorm.DB
}

type UserDetailRepository interface {
	GetUserInfoDetailInfo(ctx context.Context, entity entity.GetUserDetailInfoEntity) ([]entity.UserDetailInfoEntity, error)
}
