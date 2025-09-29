package loader

import (
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"

	"gorm.io/gorm"
)

type AuthTokenLoader struct {
	db      *gorm.DB
	storage storage.AuthTokenStorage
}

func NewAuthTokenLoader(db *gorm.DB, storage storage.AuthTokenStorage) *AuthTokenLoader {
	return &AuthTokenLoader{
		db:      db,
		storage: storage,
	}
}

func (l *AuthTokenLoader) Load(ctx context.Context) error {
	repo := repository.NewTokenRepository(l.db)

	entities, err := repo.InitUserAuthToken()
	if err != nil {
		return err
	}

	for _, v := range entities {
		log.Printf("auth token loader id : %s, uuid : %s \n", v.Id, v.Uuid)
		l.storage.PutAccessToken(v.Id, v.Uuid, v.At)
		l.storage.PutRefreshToken(v.Id, v.Uuid, v.Rt)
	}
	return nil
}
