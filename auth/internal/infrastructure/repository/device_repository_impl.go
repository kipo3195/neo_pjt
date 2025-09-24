package repository

import (
	"auth/internal/consts"
	"auth/internal/domain/device/entity"
	"auth/internal/domain/device/repository"
	"auth/internal/infrastructure/model"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) repository.DeviceRepository {
	return &deviceRepository{
		db: db,
	}
}

func (r *deviceRepository) CheckDeviceRegist(ctx context.Context, entity entity.DeviceRegistStateEntity) (bool, error) {

	var deviceRegistHis model.DeviceRegistHistory

	log.Println("디바이스 등록 조회 id :", entity.Id)

	result := r.db.Where("uuid = ? and id = ?", entity.Uuid, entity.Id).First(&deviceRegistHis)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[CheckDeviceRegist] result record = 0")
		return false, result.Error
	} else if result.Error != nil {
		log.Println("[CheckDeviceRegist] DB error")
		return false, result.Error
	} else {
		log.Println("[CheckDeviceRegist] 등록된 device")
		return true, nil
	}
}

func (r *deviceRepository) PutDevice(ctx context.Context, entity entity.DeviceRegistEntity) error {
	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&model.DeviceRegistHistory{
		Id:        entity.Id,
		Uuid:      entity.Uuid,
		ModelName: entity.ModelName,
		Version:   entity.Version,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutDevice] - Commit failed")
		return consts.ErrDB
	}
	log.Println("[PutDevice] - Commit Success")
	return nil
}
