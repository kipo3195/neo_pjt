package repository

import (
	"context"
	"user/internal/domain/userDetail/entity"
	"user/internal/domain/userDetail/repository"
	"user/internal/infrastructure/model"

	"gorm.io/gorm"
)

type userDetailRepositoryImpl struct {
	db *gorm.DB
}

func NewUserDetailRepository(db *gorm.DB) repository.UserDetailRepository {
	return &userDetailRepositoryImpl{
		db: db,
	}
}

func UserDetailMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.UserDetail{})
}

func (r *userDetailRepositoryImpl) GetUserInfoDetailInfo(ctx context.Context, en entity.GetUserDetailInfoEntity) ([]entity.UserDetailInfoEntity, error) {

	var userDetailModel []model.UserDetail

	// DB 조회
	if err := r.db.
		Table("user_detail AS a").
		Joins("JOIN service_users AS b ON a.user_hash = b.user_hash").
		Where("b.user_hash IN ?", en.UserHashs).
		Find(&userDetailModel).Error; err != nil {
		return nil, err
	}

	// 모델 → 엔티티 변환
	var userDetailEntities []entity.UserDetailInfoEntity
	for _, m := range userDetailModel {
		userDetailEntities = append(userDetailEntities, entity.UserDetailInfoEntity{
			UserHash:     m.UserHash,
			UserEmail:    m.UserEmail,
			UserPhoneNum: m.UserPhoneNum,
			// 필요한 필드만 매핑
			// 많아지면 github.com/jinzhu/copier 고려?
		})
	}

	return userDetailEntities, nil
}
