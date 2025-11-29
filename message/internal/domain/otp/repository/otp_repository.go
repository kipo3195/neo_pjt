package repository

import (
	"context"
	"message/internal/domain/otp/entity"
)

type OtpRepository interface {
	SaveOtpKey(ctx context.Context, entity entity.OTPKeyRegistEntity) error
	GetMyOtpInfo(ctx context.Context, en entity.MyOtpInfoEntity, kind string, svVersion string) (entity.OtpKeyInfoEntity, error)
}
