package repository

import (
	"auth/internal/domain/userAuth/entity"
	"context"

	"gorm.io/gorm"
)

type userAuthRepository struct {
	db *gorm.DB
}

type UserAuthRepository interface {
	PutUserAuthInfo(ctx context.Context, entity entity.UserAuthInfoEntity) error
	GetUserSalt(ctx context.Context, Id string) (string, error)
	GetUserAuthHash(ctx context.Context, Id string) (string, error)
}
