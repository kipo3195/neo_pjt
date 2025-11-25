package repository

import (
	"auth/internal/domain/otp/entity"
	"context"
)

type OtpApiRepository interface {
	OtpKeyRegistInMessage(ctx context.Context, entity entity.OtpKeyRegistEntity) (entity.OtpKeyRegistResultEntity, error)
}
