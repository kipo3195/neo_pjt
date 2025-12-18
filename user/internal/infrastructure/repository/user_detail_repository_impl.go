package repository

import (
	"context"
	"log"
	"user/internal/domain/userDetail/entity"
	"user/internal/domain/userDetail/repository"
	"user/internal/infrastructure/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	var userDetailEntities []entity.UserDetailInfoEntity

	// DB 조회
	if err := r.db.
		Table("user_detail AS a").
		Joins("JOIN service_users AS b ON a.user_hash = b.user_hash").
		Where("b.user_hash IN ?", en.UserHashs).
		Scan(&userDetailEntities).Error; err != nil {
		return nil, err
	}

	log.Println("여기 ")
	// 모델 → 엔티티 변환
	// for _, m := range userDetailEntities {
	// 	userDetailEntities = append(userDetailEntities, entity.UserDetailInfoEntity{
	// 		UserHash:     m.UserHash,
	// 		UserEmail:    m.UserEmail,
	// 		UserPhoneNum: m.UserPhoneNum,
	// 		// 필요한 필드만 매핑
	// 		// 많아지면 github.com/jinzhu/copier 고려?
	// 	})
	// }

	return userDetailEntities, nil
}

func (r *userDetailRepositoryImpl) RegistUserDetail(ctx context.Context, entity []entity.UserDetailBatchEntity) error {

	if len(entity) == 0 {
		return nil
	}

	models := make([]model.UserDetail, 0, len(entity))

	for _, e := range entity {
		models = append(models, model.UserDetail{
			Org:          e.Org,
			UserHash:     e.UserHash,
			UserPhoneNum: e.UserPhoneNum,
			UserEmail:    e.UserEmail,
		})
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "org"},       // PK 기준
				{Name: "user_hash"}, // PK 기준
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"user_phone_num",
				"user_email",
			}),
		}).
		Create(&models).Error

	if err != nil {
		return err
	}

	return nil
}
