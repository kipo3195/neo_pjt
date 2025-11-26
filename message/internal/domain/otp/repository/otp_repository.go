package repository

import (
	"context"
	"message/internal/domain/otp/entity"
)

type OtpRepository interface {
	SaveOtpKey(ctx context.Context, entity *entity.OTPKeyRegistEntity) error
}
