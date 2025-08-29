package repository

import (
	"context"
	"org/internal/domain/department/entities"
)

type DepartmentRepository interface {
	PutDept(ctx context.Context, entity entities.CreateDeptEntity) (interface{}, error)
	DeleteDept(ctx context.Context, entity entities.DeleteDeptEntity) (interface{}, error)
	PutDeptUser(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error)
	DeleteDeptUser(ctx context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error)
}
