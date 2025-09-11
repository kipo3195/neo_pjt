package repository

import (
	"common/internal/domain/user/entity"
	"context"
)

type UserAPIRepository interface {
	UserAuthRegistInAuth(ctx context.Context, id string, entity entity.UserRegisterInfoEntity, challenge string) (string, error)
}
