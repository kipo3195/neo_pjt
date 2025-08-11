package loader

import (
	repository "common/internal/domains/skin/repositories/server"
	"common/internal/infra/storage"
	"context"

	"gorm.io/gorm"
)

type SkinLoader struct {
	db      *gorm.DB
	storage *storage.SkinStorage
}

func NewSkinLoader(db *gorm.DB, storage *storage.SkinStorage) *SkinLoader {
	return &SkinLoader{db: db, storage: storage}
}

func (l *SkinLoader) Load(ctx context.Context) error {
	skins, err := repository.SaveConfigHash(l.db) // ← 도메인별 repository
	if err != nil {
		return err
	}
	l.storage.SaveConfigHash(skins)
	return nil
}
