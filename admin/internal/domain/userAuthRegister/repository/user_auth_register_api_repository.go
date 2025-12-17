package repository

import (
	"context"
)

type UserAuthRegisterApiRepository interface {
	UserAuthRegisterInAuth(ctx context.Context, entity []entity.UserAuthRegisterEntity) error
}
