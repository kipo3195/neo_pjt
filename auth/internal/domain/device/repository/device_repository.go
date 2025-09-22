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
