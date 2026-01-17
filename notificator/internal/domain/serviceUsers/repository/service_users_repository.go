package repository

import (
	"context"
	"notificator/internal/domain/serviceUsers/entity"
)

type ServiceUsersRepository interface {
	PutServiceUser(ctx context.Context, en []entity.RegisterServiceUsersEntity) error
}
