package repository

import (
	"batch/internal/domain/userDetail/entity"
	"context"
)

type UserDetailRepository interface {
	GetUserDetail(ctx context.Context, org string) ([]entity.UserDetailEntity, error)
}
