package loader

import (
	"auth/internal/infrastructure/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"

	"gorm.io/gorm"
)

type ServiceUserLoader struct {
	db      *gorm.DB
	storage storage.ServiceUserStorage
}

func NewServiceUserLoader(db *gorm.DB, storage storage.ServiceUserStorage) *ServiceUserLoader {

	return &ServiceUserLoader{
		db:      db,
		storage: storage,
	}
}

func (l *ServiceUserLoader) Load(ctx context.Context) error {
	repo := repository.NewServiceUserRepository(l.db)

	serviceUsers, err := repo.InitServiceUsers()
	if err != nil {
		return err
	}
	log.Printf("auth serviceUserLoader loader start. \n")
	for _, v := range serviceUsers {
		l.storage.PutUserhash(v.UserId, v.UserHash)
	}
	log.Printf("auth serviceUserLoader loader end. count : %d \n", len(serviceUsers))
	return nil
}
