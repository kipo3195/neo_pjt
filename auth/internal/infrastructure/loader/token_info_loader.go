package loader

import (
	"auth/internal/consts"
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"

	"gorm.io/gorm"
)

type TokenInfoLoader struct {
	db      *gorm.DB
	storage storage.AuthTokenStorage
}

func NewTokenInfoLoader(db *gorm.DB, storage storage.AuthTokenStorage) *TokenInfoLoader {
	return &TokenInfoLoader{
		db:      db,
		storage: storage,
	}
}

func (l *TokenInfoLoader) Load(ctx context.Context) error {

	repo := repository.NewTokenRepository(l.db)
	entities, err := repo.InitAuthTokenInfo(ctx)

	if err != nil {
		// 에러
		// 아래 로직이 주석처리된 이유 : l.storage에서 GetTokenExpInfo 로직에서 default value를 init하고 있으므로..
		// log.Println("[InitAuthTokenInfo] DB error. default value init.")
		// l.storage.PutTokenExpInfo(consts.DEVICE_ACCESSS_TOKEN, 60)
		// l.storage.PutTokenExpInfo(consts.DEVICE_REFRESH_TOKEN, 30)
		return err
	}

	for _, v := range entities {
		log.Println("TokenType : ", v.TokenType, "TokenExp : ", v.TokenExp)
		if v.TokenType == consts.DEVICE_ACCESSS_TOKEN {
			l.storage.PutTokenExpInfo(consts.DEVICE_ACCESSS_TOKEN, v.TokenExp)
		} else if v.TokenType == consts.DEVICE_REFRESH_TOKEN {
			l.storage.PutTokenExpInfo(consts.DEVICE_REFRESH_TOKEN, v.TokenExp)
		}
	}

	return nil
}
