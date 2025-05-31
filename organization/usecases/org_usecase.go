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
	ServerDeleteDepartment(ctx context.Context, req svDto.ServerDeleteDeptRequest) (interface{}, error)
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) ServerCreateDepartment(ctx context.Context, req svDto.ServerCreateDeptRequest) (interface{}, error) {
	return r.repo.SaveDepartment(ctx, ToCreateDepartmentEntity(req))
}

func ToCreateDepartmentEntity(req svDto.ServerCreateDeptRequest) entities.CreateDepartmentEntity {

	return entities.CreateDepartmentEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		DeptNameKr:     req.DeptNameKr,
		DeptNameEn:     req.DeptNameEn,
		DeptNameJp:     req.DeptNameJp,
		DeptNameCn:     req.DeptNameCn,
	}
}

func (r *orgUsecase) ServerDeleteDepartment(ctx context.Context, req svDto.ServerDeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDepartment(ctx, ToDeleteDepartmentEntity(req))
}

func ToDeleteDepartmentEntity(req svDto.ServerDeleteDeptRequest) entities.DeleteDepartmentEntity {

	return entities.DeleteDepartmentEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}
