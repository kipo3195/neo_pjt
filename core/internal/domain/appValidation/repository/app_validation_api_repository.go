package repository

import (
	"context"
	"core/internal/domain/appValidation/entity"
)

type AppValidationAPIRepository interface {
	DeviceInit(ctx context.Context, entity entity.ValidationEntity) (*entity.DeviceInitResult, error)
}
