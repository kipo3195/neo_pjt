package repository

import (
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
)

type TokenRepository interface {
	PutIssuedAppToken(token *shared.AppTokenEntity) (bool, error)
	GetValidationAppToken(entity entity.AppTokenValidationEntity) (bool, error)
}
