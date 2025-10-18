package repository

import (
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
	"context"
)

type TokenRepository interface {
	PutIssuedAppToken(token *shared.AppTokenEntity) (bool, error)
	GetValidationAppToken(entity entity.AppTokenValidationEntity) (bool, error)
	InitUserAuthToken() ([]entity.AuthTokenEntity, error)
	InitAuthTokenInfo(ctx context.Context) ([]entity.AuthTokenInfoEntity, error)
	PutAuthToken(ctx context.Context, id string, uuid string, at string, rt string, rtExp string) error
	UpdateReIssueAccessTokenInfo(ctx context.Context, entity entity.ReIssueAccessTokenSavedEntity) error
	GetUserIdWithRtAndUuid(ctx context.Context, entity entity.RefreshTokenCheckEntity) (string, error)
}
