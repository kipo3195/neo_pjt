package repository

import (
	"auth/internal/domain/device/entity"
	"context"
)

type DeviceRepository interface {
	CheckDeviceRegist(ctx context.Context, entity entity.DeviceRegistStateEntity) (bool, error)
	PutDevice(ctx context.Context, entity entity.DeviceRegistEntity) error
}
