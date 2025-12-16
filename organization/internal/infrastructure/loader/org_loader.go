package loader

import (
	"context"
	"org/internal/infrastructure/repository"
	"org/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type OrgLoader struct {
	db      *gorm.DB
	storage storage.OrgStorage
}

func NewOrgLoader(db *gorm.DB, storage storage.OrgStorage) *OrgLoader {

	return &OrgLoader{
		db:      db,
		storage: storage,
	}
}

func (l *OrgLoader) Load(ctx context.Context) error {

	repo := repository.NewOrgRepository(l.db)
	worksOrgCode, err := repo.InitWorksOrgCode(ctx)
	if err != nil {
		return err
	}

	l.storage.PutWorksOrgCode(worksOrgCode)
	return nil
}
