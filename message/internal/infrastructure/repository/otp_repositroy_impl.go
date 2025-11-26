package repository

import (
	"message/internal/domain/otp/repository"

	"gorm.io/gorm"
)

type otpApiRepositoryImpl struct {
	db *gorm.DB
}

func NewOtpApiRepository(db *gorm.DB) repository.OtpRepository {
	return &otpApiRepositoryImpl{
		db: db,
	}
}
