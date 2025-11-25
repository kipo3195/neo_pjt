package repository

import (
	"auth/internal/domain/otp/entity"
	"auth/internal/domain/otp/repository"
	"context"
)

type otpApiRepositoryImpl struct {
}

func NewOtpApiRepository() repository.OtpApiRepository {
	return &otpApiRepositoryImpl{}
}

func (r *otpApiRepositoryImpl) OtpKeyRegistInMessage(ctx context.Context, en entity.OtpKeyRegistEntity) (entity.OtpKeyRegistResultEntity, error) {

	// message API 호출

	result := entity.OtpKeyRegistResultEntity{}

	return result, nil
}
