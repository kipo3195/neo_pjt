package loader

import (
	"common/internal/consts"
	repository "common/internal/domains/configuration/repositories/client"
	"common/internal/infra/storage"
	"context"

	"gorm.io/gorm"
)

type ConfigHashLoader struct {
	db      *gorm.DB
	storage storage.ConfigHashStorage
}

func NewConfigHashLoader(db *gorm.DB, storage storage.ConfigHashStorage) *ConfigHashLoader {
	return &ConfigHashLoader{db: db, storage: storage}
}

func (l *ConfigHashLoader) Load(ctx context.Context) error {

	// DDD나 클린 아키텍처 흐름에서 Repository는 보통 Usecase 레이어에 주입되지만,
	// Loader는 서버 부팅 시점의 "시스템 초기화" 로직이기 때문에
	// Handler도 없고 Usecase도 호출하지 않고
	// 바로 DB → Storage로 가는 단방향 로직만 필요합니다.
	// 그래서 Loader가 직접 repository를 생성하거나,
	// 아예 repository가 아니라 DAO(Data Access Object)처럼 DB에서 가져오는 함수만 호출하는 식으로 단순화하는 경우가 많습니다.

	repo := repository.NewConfigurationRepository(l.db)
	configs, err := repo.GetConfigHash()
	if err != nil {
		return err
	}
	l.storage.SaveConfigHash(consts.CONFIG, configs)
	return nil
}
