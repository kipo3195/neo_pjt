package usecases

import (
	"admin/dto"
	"admin/repositories"
	"context"
)

type adminUsecase struct {
	repo repositories.AdminRepository
}

type AdminUsecase interface {
	CreateDepartment(ctx context.Context, req dto.CreateDeptRequest) (interface{}, error)
}

func NewAdminUsecase(repo repositories.AdminRepository) AdminUsecase {
	return &adminUsecase{repo: repo}
}

func (r *adminUsecase) CreateDepartment(ctx context.Context, req dto.CreateDeptRequest) (interface{}, error) {

	return nil, nil
}
