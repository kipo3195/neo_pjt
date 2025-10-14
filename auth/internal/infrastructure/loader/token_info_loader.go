package loader

import (
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"

	"gorm.io/gorm"
)

type TokenInfoLoader struct {
	db      *gorm.DB
	storage storage.TokenStorage
}

func NewTokenInfoLoader(db *gorm.DB, storage storage.TokenStorage) *TokenInfoLoader {
	return &TokenInfoLoader{
		db:      db,
		storage: storage,
	}
}

func (l *TokenInfoLoader) Load(ctx context.Context) error {
	repo := repository.NewTokenRepository(l.db)

	entities, err := repo.InitAuthTokenInfo(ctx)
	if err != nil {
		return err
	}

	for _, v := range entities {
		log.Println("TokenType : ", v.TokenType, "TokenExp : ", v.TokenExp)
		l.storage.PutTokenExp(v.TokenType, v.TokenExp)
	}

	return nil
}
