package repository

import (
	"context"
	"log"
	"message/internal/consts"
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

	modelEntity := model.OtpKeyRegistHistory{
		Id:           entity.Id,
		Uuid:         entity.Uuid,
		ChatOtpKey:   entity.ChatOtpKey,
		NoteOtpKey:   entity.NoteOtpKey,
		RegDate:      entity.OtpRegDate,
		SvKeyVersion: entity.SvKeyVersion,
	}

	// Insert, 중복(PK 충돌) 시 ChatOtpKey, NoteOtpKey, RegDate만 update
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}, {Name: "uuid"}, {Name: "sv_key_version"}}, // UNIQUE KEY or PRIMARY KEY 기준
		DoUpdates: clause.Assignments(map[string]interface{}{
			"chat_otp_key": entity.ChatOtpKey,
			"note_otp_key": entity.NoteOtpKey,
			"reg_date":     entity.OtpRegDate,
		}),
	}).Create(&modelEntity).Error

	return err
}

func (r *otpApiRepositoryImpl) GetMyOtpInfoLatest(ctx context.Context, en entity.MyOtpInfoEntity, svVersion string) ([]entity.MyOtpInfoResultEntity, error) {
	var modelEntities []model.OtpKeyRegistHistory
	err := r.db.WithContext(ctx).Where("id = ? and uuid = ? and sv_key_version = ?", en.UserId, en.Uuid, svVersion).Find(&modelEntities).Error
	if err != nil {
		return nil, err
	}
	if len(modelEntities) == 0 {
		log.Println("GetMyOtpInfoLatest: no data found")
		return nil, consts.ErrDBresultNotFound
	}

	modelEntity := modelEntities[0]
	result := make([]entity.MyOtpInfoResultEntity, 0)

	result = append(result, entity.MyOtpInfoResultEntity{
		Version:    modelEntity.SvKeyVersion,
		KeyType:    consts.CHAT,
		Key:        modelEntity.ChatOtpKey,
		OtpRegDate: modelEntity.RegDate,
	})

	result = append(result, entity.MyOtpInfoResultEntity{
		Version:    modelEntity.SvKeyVersion,
		KeyType:    consts.NOTE,
		Key:        modelEntity.NoteOtpKey,
		OtpRegDate: modelEntity.RegDate,
	})

	return result, nil
}
