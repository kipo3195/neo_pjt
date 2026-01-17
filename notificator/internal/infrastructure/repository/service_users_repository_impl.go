package repository

import (
	"context"
	"log"
	"notificator/internal/domain/serviceUsers/entity"
	"notificator/internal/domain/serviceUsers/repository"
	"notificator/internal/infrastructure/model"

	"gorm.io/gorm"
)

type serviceUsersRepositoryImpl struct {
	db *gorm.DB
}

func ServiceUsersMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ServiceUsers{})
}

func NewServiceUsersRepository(db *gorm.DB) repository.ServiceUsersRepository {
	return &serviceUsersRepositoryImpl{
		db: db,
	}
}

func (r *serviceUsersRepositoryImpl) PutServiceUser(ctx context.Context, en []entity.RegisterServiceUsersEntity) error {

	var models []model.ServiceUsers

	for _, e := range en {
		models = append(models, model.ServiceUsers{
			Org:      e.Org,
			UserId:   e.UserId,
			UserHash: e.UserHash,
		})
	}

	// insert + update로 처리하지 않는이유
	// auto scale된 상황에서 admin이 2개 이상떴을때 특정 userId추가 됬음에도 service_users 조회시 실시간 반영되지 않은 상황에서 추가 처리 X를 위함.
	// 그래야 NATS를 통해 데이터를 보낼때는 성공한 CASE에 대해서만 보낼 수 있으므로.

	if err := r.db.Create(&models).Error; err != nil {
		log.Println("[PutServiceUser] err : ", err)
		return err
	}

	return nil
}
