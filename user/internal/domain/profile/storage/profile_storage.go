package storage

import (
	"context"
	"user/internal/domain/profile/entity"
)

type ProfileStorage interface {
	Upload(ctx context.Context, profileImg []byte, fileHash string) (string, string, error)
	GetProfileUrl(ctx context.Context, fileName string) ([]byte, error)
	DeleteImg(ctx context.Context, fileName string) error
	GetUserProfileUpdateHash(ctx context.Context, entity []entity.ReqUserEntity) (map[string]entity.ProfileHashEntity, error)
	GetProfileInfo(userHash []string) []entity.ProfileInfoEntity
}
