package repository

import (
	"common/internal/domain/skin/entity"
	"context"

	"gorm.io/gorm"
)

type skinRepository struct {
	db *gorm.DB
}

type SkinRepository interface {
	GetSkinHash() (string, error)
	PutSkinFileInfo(ctx context.Context, entity *entity.SkinFileInfoEntity) (bool, error)
}
