package loader

import (
	"common/internal/infrastructure/repository"
	"common/internal/infrastructure/storage"
	"context"

	"gorm.io/gorm"
)

type OrgLoader struct {
	db      *gorm.DB
	storage storage.OrgStorage
}

func NewOrgLoader(db *gorm.DB, storage storage.OrgStorage) *OrgLoader {
	return &OrgLoader{db: db, storage: storage}
}

func (l *OrgLoader) Load(ctx context.Context) error {

	repo := repository.NewOrgRepository(l.db)
	orgCode, err := repo.GetOrgCode()
	if err != nil {
		return err
	}

	l.storage.PutOrgInfo(orgCode)
	return nil
}
