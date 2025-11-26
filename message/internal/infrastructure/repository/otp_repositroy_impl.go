package repository

import (
	"context"
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

func (r *otpApiRepositoryImpl) SaveOtpKey(ctx context.Context, entity *entity.OTPKeyRegistEntity) error {

	modelEntity := model.OtpKeyRegistHistory{
		Id:         entity.Id,
		Uuid:       entity.Uuid,
		ChatOtpKey: entity.ChatOtpKey,
		NoteOtpKey: entity.NoteOtpKey,
		RegDate:    entity.RegDate,
	}

	// Insert, 중복(PK 충돌) 시 ChatOtpKey, NoteOtpKey, RegDate만 update
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}, {Name: "uuid"}}, // UNIQUE KEY or PRIMARY KEY 기준
		DoUpdates: clause.Assignments(map[string]interface{}{
			"chat_otp_key": entity.ChatOtpKey,
			"note_otp_key": entity.NoteOtpKey,
			"reg_date":     entity.RegDate,
		}),
	}).Create(&modelEntity).Error

	return err
}
