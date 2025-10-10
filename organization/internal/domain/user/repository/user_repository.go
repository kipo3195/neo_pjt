package repository

import (
	"context"
	"org/internal/domain/user/entity"
)

type UserRepository interface {
	GetMyInfo(ctx context.Context, entity entity.MyInfoHashEntity) (entity.MyInfoEntity, error)
	CreateServiceUser(ctx context.Context, entity []entity.ServiceUserEntity) error
	GetServiceUsers(ctx context.Context, keyword string) ([]entity.UserDetailEntity, error)
	CreateUserDetail(ctx context.Context, entity []entity.UserDetailEntity) error
	CreateUserMultiLang(ctx context.Context, entity entity.UserMultiLangEntity) error
}
