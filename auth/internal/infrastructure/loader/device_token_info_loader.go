package loader

import (
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"

	"gorm.io/gorm"
)

type DeviceTokenInfoLoader struct {
	db      *gorm.DB
	storage storage.DeviceStorage
}

func NewDeviceTokenInfoLoader(db *gorm.DB, storage storage.DeviceStorage) *DeviceTokenInfoLoader {
	return &DeviceTokenInfoLoader{
		db:      db,
		storage: storage,
	}
}

func (l *DeviceTokenInfoLoader) Load(ctx context.Context) error {
	repo := repository.NewDeviceRepository(l.db)

	entities, err := repo.InitDeviceTokenInfo(ctx)
	if err != nil {
		return err
	}

	for _, v := range entities {
		log.Println("TokenType : ", v.TokenType, "TokenExp : ", v.TokenExp)
		l.storage.PutDeviceTokenExp(v.TokenType, v.TokenExp)
	}

	return nil
}
