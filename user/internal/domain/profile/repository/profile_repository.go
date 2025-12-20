package repository

import (
	"context"
	"user/internal/domain/profile/entity"

	"gorm.io/gorm"
)

type profileRepositroy struct {
	db *gorm.DB
}

type ProfileRepository interface {
	PutUserProfileImgInfo(ctx context.Context, entity entity.ProfileImgEntity) error
	DeleteUserProfileImgInfo(ctx context.Context, userHash string, fileName string) error
	RollbackDeleteUserProfileImgInfo(ctx context.Context, userHash string, fileName string) error
	GetProfileInfo(ctx context.Context, entity entity.GetProfileInfoEntity) (map[string]entity.GetProfileInfoResultEntity, error)
	PutProfileMsg(ctx context.Context, entity entity.PutProfileMsgEntity) error
	GetProfileMsg(ctx context.Context, entity entity.GetProfileMsgEntity) ([]entity.GetProfileMsgResultEntity, error)
	InitProfile(ctx context.Context) ([]entity.ProfileInfoEntity, error)
}
