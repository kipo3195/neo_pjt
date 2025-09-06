package repository

import (
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"

	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

type TokenRepository interface {
	PutIssuedAppToken(token *shared.AppTokenEntity) (bool, error)
	GetValidationAppToken(entity entity.AppTokenValidationEntity) (bool, error)
}
