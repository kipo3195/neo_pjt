package loader

import (
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"
	"context"

	"gorm.io/gorm"
)

type SkinLoader struct {
	db      *gorm.DB
	storage storage.SkinStorage
}

func NewSkinLoader(db *gorm.DB, storage storage.SkinStorage) *SkinLoader {
	return &SkinLoader{db: db, storage: storage}
}

func (l *SkinLoader) Load(ctx context.Context) error {

	repo := repository.NewSkinRepository(l.db)

	skinHash, err := repo.GetSkinHash() // ← 도메인별 repository
	if err != nil {
		return err

	}

	l.storage.SaveSkinHash("skin", skinHash)
	return nil
}
