package repository

import (
	"context"
	"log"
	"message/internal/domain/otp/entity"
	"message/internal/domain/otp/repository"
	"message/internal/infrastructure/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type otpApiRepositoryImpl struct {
	db *gorm.DB
}

func OtpKeyMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.OtpKeyRegistHistory{})
}

func NewOtpApiRepository(db *gorm.DB) repository.OtpRepository {
	return &otpApiRepositoryImpl{
		db: db,
	}
}

func (r *otpApiRepositoryImpl) SaveOtpKey(ctx context.Context, entity entity.OTPKeyRegistEntity) error {

	log.Println("[SaveOtpKey] DB insert")
	id := entity.Id
	uuid := entity.Uuid

	for _, info := range entity.OtpKeyInfoEntity {

		// DB row 구조체로 변환
		row := model.OtpKeyRegistHistory{
			Id:           id,
			Uuid:         uuid,
			OtpKey:       info.OtpKey,
			Kind:         info.Kind,
			RegDate:      info.OtpRegDate,
			SvKeyVersion: info.SvKeyVersion,
		}

		// Upsert 실행
		err := r.db.Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "id"},
					{Name: "uuid"},
					{Name: "sv_key_version"},
					{Name: "kind"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"otp_key",
					"reg_date",
				}),
			},
		).Create(&row).Error

		if err != nil {
			return err
		}
		log.Printf("[SaveOtpKey] id :%s, uuid :%s, kind :%s insert success. \n ", id, uuid, info.Kind)
	}

	return nil
}

func (r *otpApiRepositoryImpl) GetMyOtpInfo(ctx context.Context, en entity.MyOtpInfoEntity, kind string, svVersion string) (entity.OtpKeyInfoEntity, error) {
	var modelEntity model.OtpKeyRegistHistory

	err := r.db.WithContext(ctx).Where("id = ? and uuid = ? and kind = ? and sv_key_version = ?", en.UserId, en.Uuid, kind, svVersion).Take(&modelEntity).Error

	if err != nil {
		return entity.OtpKeyInfoEntity{}, err
	}

	result := entity.OtpKeyInfoEntity{
		OtpKey:       modelEntity.OtpKey,
		Kind:         modelEntity.Kind,
		OtpRegDate:   modelEntity.RegDate,
		SvKeyVersion: modelEntity.SvKeyVersion,
	}

	return result, nil
}
