package repository

import (
	"admin/internal/domain/serviceUser/entity"
	"context"
)

type ServiceUserRepository interface {
	PutServiceUser(ctx context.Context, org string, entity []entity.ServiceUserEntity) ([]entity.ServiceUserEntity, error)
}
