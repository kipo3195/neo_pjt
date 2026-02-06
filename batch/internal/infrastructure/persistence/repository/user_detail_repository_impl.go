package repository

import (
	"batch/internal/domain/userDetail/entity"
	"batch/internal/domain/userDetail/repository"
	"batch/internal/infrastructure/persistence/model"
	"context"
	"log"

	"gorm.io/gorm"
)

type userDetailRepositoryImpl struct {
	db *gorm.DB
}

func UserDetailMigrate(db *gorm.DB) error {
	return db.AutoMigrate(model.UserDetail{})
}

func NewUserDetailRepository(db *gorm.DB) repository.UserDetailRepository {
	return &userDetailRepositoryImpl{
		db: db,
	}
}

func (r *userDetailRepositoryImpl) GetUserDetail(ctx context.Context, org string) ([]entity.UserDetailEntity, error) {

	userDetailModel := []model.UserDetail{}

	viewSql := `select * from user_detail_view where org = ?`
	err := r.db.Raw(viewSql, org).Scan(&userDetailModel).Error

	if err != nil {
		log.Println("[GetUserDetail] - No record found or DB error")
		return nil, err
	}

	en := make([]entity.UserDetailEntity, 0)

	for _, m := range userDetailModel {

		temp := entity.UserDetailEntity{
			Org:          m.Org,
			UserId:       m.UserId,
			UserHash:     m.UserHash,
			UserEmail:    m.UserEmail,
			UserPhoneNum: m.UserPhoneNum,
		}

		en = append(en, temp)
	}

	return en, nil
}
