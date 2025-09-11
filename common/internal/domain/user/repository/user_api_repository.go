package repository

import (
	"common/internal/domain/user/entity"
	"context"
)

type UserAPIRepository interface {
	UserRegistInAuth(ctx context.Context, id string, entity entity.UserRegisterInfoEntity) (string, error)
}
