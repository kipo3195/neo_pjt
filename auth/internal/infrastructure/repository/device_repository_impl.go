package repository

import (
	"auth/internal/domain/device/entity"

	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

type DeviceRepository interface {
	CheckDeviceRegist(entity entity.DeviceEntity) (bool, error)
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{
		db: db,
	}
}

func (r *deviceRepository) CheckDeviceRegist(entity entity.DeviceEntity) (bool, error) {

	return false, nil
}
