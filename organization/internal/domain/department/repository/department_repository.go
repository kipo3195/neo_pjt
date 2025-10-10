package repository

import (
	"context"
	"org/internal/domain/department/entity"
)

type DepartmentRepository interface {
	PutDept(ctx context.Context, entity entity.CreateDeptEntity) (interface{}, error)
	DeleteDept(ctx context.Context, entity entity.DeleteDeptEntity) (interface{}, error)
	PutDeptUser(ctx context.Context, entity entity.CreateDeptUserEntity) (interface{}, error)
	DeleteDeptUser(ctx context.Context, entity entity.DeleteDeptUserEntity) (interface{}, error)
	CreateDeptTree(ctx context.Context, entity entity.WorksDeptEntity) error
	GetDepts(ctx context.Context, org string) ([]entity.WorksDeptEntity, error)
	PutWorksDeptMultiLang(ctx context.Context, entity entity.CreateMultiLangEntity) error
}
