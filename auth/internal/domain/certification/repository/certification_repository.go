package repository

import (
	"auth/internal/domain/certification/entity"

	"gorm.io/gorm"
)

type cerificationRepository struct {
	db *gorm.DB
}

type CerificationRepository interface {
	CheckAuth(entity entity.AuthInfoEntity) (*entity.AuthInfoEntity, error)
	GetUserHash(entity entity.AuthInfoEntity) (string, error)
	GetValidation(entity entity.AppTokenValidationEntity) (bool, error)
}
