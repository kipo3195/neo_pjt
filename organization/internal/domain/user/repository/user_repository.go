package repository

import (
	"context"
	sharedEntity "org/internal/domain/shared/entity"
)

type UserRepository interface {
	GetMyInfo(ctx context.Context, myHash string) (sharedEntity.MyInfoEntity, error)
}
