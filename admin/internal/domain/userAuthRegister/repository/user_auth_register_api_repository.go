package repository

import (
	"admin/internal/domain/userAuthRegister/entity"
	"context"
)

type UserAuthRegisterApiRepository interface {
	UserAuthRegisterInAuth(ctx context.Context, entity []entity.UserAuthRegisterEntity) (string, error)
}
