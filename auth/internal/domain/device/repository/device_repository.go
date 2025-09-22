package repository

import (
	"auth/internal/domain/device/entity"
)

type DeviceRepository interface {
	CheckDeviceRegist(entity entity.DeviceEntity) (bool, error)
}
