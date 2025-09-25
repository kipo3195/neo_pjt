package repository

import (
	"context"
	"org/internal/domain/user/entity"
)

type UserRepository interface {
	GetMyInfo(ctx context.Context, entity entity.MyInfoHashEntity) (entity.MyInfoEntity, error)
}
