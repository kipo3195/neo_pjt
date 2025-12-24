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
	RegistUserDetail(ctx context.Context, entity []entity.UserDetailBatchEntity) error
	InitUserDetail(ctx context.Context) ([]entity.UserDetailInfoEntity, error)
	GetOrgUsers(ctx context.Context, orgCode string) ([]entity.UserDetailInfoEntity, error)
}
