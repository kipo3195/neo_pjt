package repository

import (
	"auth/internal/domain/serviceUser/entity"
	"auth/internal/domain/serviceUser/repository"
	"auth/internal/infrastructure/model"

	"gorm.io/gorm"
)

type serviceUserRepository struct {
	db *gorm.DB
}

func NewServiceUserRepository(db *gorm.DB) repository.ServiceUserRepository {
	return &serviceUserRepository{
		db: db,
	}
}

func (r *serviceUserRepository) InitServiceUsers() ([]entity.ServiceUserEntity, error) {

	// 한명의 사용자가 여러개의 ID를 갖더라도 hash는 하나로..
	// seq를 pk로 갖고 id : hash 형태를 갖는 테이블이면 될듯
	var users []model.ServiceUsers

	// 전체 조회
	if err := r.db.Where(`use_yn = ?`, "Y").Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]entity.ServiceUserEntity, 0, len(users))
	for _, u := range users {
		result = append(result, entity.ServiceUserEntity{
			UserHash: u.UserHash,
			UserId:   u.UserId,
		})
	}

	return result, nil
}
