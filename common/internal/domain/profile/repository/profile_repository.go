package repository

import (
	"common/internal/domain/profile/entity"
	"context"

	"gorm.io/gorm"
)

type profileRepositroy struct {
	db *gorm.DB
}

type ProfileRepository interface {
	PutUserProfileImgInfo(ctx context.Context, entity entity.ProfileImgEntity) error
	DeleteUserProfileImgInfo(ctx context.Context, userId string, fileName string) error
}
