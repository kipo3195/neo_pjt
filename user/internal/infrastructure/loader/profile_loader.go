package loader

import (
	"context"
	"log"
	"time"
	"user/internal/infrastructure/repository"
	"user/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type ProfileLoader struct {
	db      *gorm.DB
	storage storage.ProfileCacheStorage
}

func NewProfileLoader(db *gorm.DB, storage storage.ProfileCacheStorage) *ProfileLoader {

	return &ProfileLoader{
		db:      db,
		storage: storage,
	}

}

func (l *ProfileLoader) Load(ctx context.Context) error {

	repo := repository.NewProfileRepository(l.db)

	entities, err := repo.InitProfile(ctx)

	if err != nil {
		log.Println("[InitProfile] init error.")
		return err
	}

	profileVersion := time.Now().UnixNano()

	l.storage.InitProfile(ctx, profileVersion, entities)

	return nil
}
