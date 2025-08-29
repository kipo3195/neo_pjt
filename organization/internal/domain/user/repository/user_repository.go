package repository

import (
	"context"
	"org/internal/sharedEntities"
)

type UserRepository interface {
	GetMyInfo(ctx context.Context, myHash string) (sharedEntities.MyInfoEntity, error)
}
