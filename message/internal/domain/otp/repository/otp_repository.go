package repository

import (
	"context"
	"message/internal/domain/otp/entity"
)

type OtpRepository interface {
	SaveOtpKey(ctx context.Context, entity entity.OTPKeyRegistEntity) error
	GetMyOtpInfoLatest(ctx context.Context, entity entity.MyOtpInfoEntity, svVersion string) (en []entity.MyOtpInfoResultEntity, err error)
}
