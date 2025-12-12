package repository

import (
	"admin/internal/domain/serviceUser/entity"
	"admin/internal/domain/serviceUser/repository"
	"admin/internal/infrastructure/model"
	"context"
	"log"

	"gorm.io/gorm"
)

type serviceUserRepositoryImpl struct {
	db *gorm.DB
}

func ServiceUsersMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ServiceUsers{})
}

func NewServiceUserRepository(db *gorm.DB) repository.ServiceUserRepository {
	return &serviceUserRepositoryImpl{
		db: db,
	}
}

func (r *serviceUserRepositoryImpl) PutServiceUser(ctx context.Context, org string, entity []entity.ServiceUserEntity) ([]entity.ServiceUserEntity, error) {

	var models []model.ServiceUsers

	for _, e := range entity {
		models = append(models, model.ServiceUsers{
			Org:      org,
			UserId:   e.UserId,
			UserHash: e.UserHash,
		})
	}

	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutServiceUser] err : ", err)
		return nil, err
	}

	return entity, nil
}
