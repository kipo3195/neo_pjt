package usecases

import (
	"context"
	svDto "org/dto/server"
	"org/entities"
	"org/repositories"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
	ServerCreateDepartment(ctx context.Context, req svDto.ServerCreateDeptRequest) (interface{}, error)
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) ServerCreateDepartment(ctx context.Context, req svDto.ServerCreateDeptRequest) (interface{}, error) {
	return r.repo.SaveDepartment(ctx, ToCreateDepartmentEntity(req))
}

func ToCreateDepartmentEntity(req svDto.ServerCreateDeptRequest) entities.CreateDepartmentEntity {

	return entities.CreateDepartmentEntity{
		DeptName: req.DeptName,
	}
}
