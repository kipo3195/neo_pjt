package loader

import (
	"context"
	"log"
	"time"
	"user/internal/infrastructure/repository"
	"user/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type UserDetailLoader struct {
	db      *gorm.DB
	storage storage.UserDetailStorage
}

func NewUserDetailLoader(db *gorm.DB, storage storage.UserDetailStorage) *UserDetailLoader {

	return &UserDetailLoader{
		db:      db,
		storage: storage,
	}
}

func (l *UserDetailLoader) Load(ctx context.Context) error {
	repo := repository.NewUserDetailRepository(l.db)

	entities, err := repo.InitUserDetail(ctx)

	log.Println("[InitUserDetail] entities size : ", len(entities))

	if err != nil {
		log.Println("[InitUserDetail] init error.")
		return err
	}

	detailVersion := time.Now().UnixNano()

	// 현재 시간 생성 - detailVersion
	l.storage.InitUserDetail(ctx, detailVersion, entities)

	return nil
}
