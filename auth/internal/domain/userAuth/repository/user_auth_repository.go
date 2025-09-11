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
	PutUserAuth(ctx context.Context, entity entity.UserAuthEntity) error
}
