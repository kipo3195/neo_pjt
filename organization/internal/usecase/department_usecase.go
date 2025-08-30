package usecase

import (
	"context"
	"log"
	"org/internal/domain/department/entity"
	"org/internal/domain/department/repository"
	"org/internal/dto/department"
	"org/internal/utils"
)

type departmentUsecase struct {
	repository repository.DepartmentRepository
}

type DepartmentUsecase interface {
	CreateDept(ctx context.Context, requestDTO department.CreateDeptRequest) (interface{}, error)

	DeleteDept(ctx context.Context, req department.DeleteDeptRequest) (interface{}, error)

	CreateDeptUser(ctx context.Context, req department.CreateDeptUserRequest) (interface{}, error)
	DeleteDeptUser(ctx context.Context, req department.DeleteDeptUserRequest) (interface{}, error)
}

func NewDepartmentUsecase(repository repository.DepartmentRepository) DepartmentUsecase {
	return &departmentUsecase{
		repository: repository,
	}
}

func (r *departmentUsecase) CreateDept(ctx context.Context, req department.CreateDeptRequest) (interface{}, error) {
	return r.repository.PutDept(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(requestDTO department.CreateDeptRequest) entity.CreateDeptEntity {

	return entity.CreateDeptEntity{
		DeptCode:       requestDTO.DeptCode,
		DeptOrg:        requestDTO.DeptOrg,
		ParentDeptCode: requestDTO.ParentDeptCode,
		KoLang:         requestDTO.KoLang,
		EnLang:         requestDTO.EnLang,
		JpLang:         requestDTO.JpLang,
		ZhLang:         requestDTO.ZhLang,
		RuLang:         requestDTO.RuLang,
		ViLang:         requestDTO.ViLang,
		Header:         requestDTO.Header,
	}
}

func (r *departmentUsecase) DeleteDeptUser(ctx context.Context, req department.DeleteDeptUserRequest) (interface{}, error) {
	return r.repository.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req department.DeleteDeptUserRequest) entity.DeleteDeptUserEntity {

	return entity.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *departmentUsecase) CreateDeptUser(ctx context.Context, req department.CreateDeptUserRequest) (interface{}, error) {

	updateHash := utils.MakeUpdateHash()
	log.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repository.PutDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req department.CreateDeptUserRequest, updateHash string) entity.CreateDeptUserEntity {

	return entity.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		DeptOrg:              req.DeptOrg,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
		UpdateHash:           updateHash,
	}
}

func (r *departmentUsecase) DeleteDept(ctx context.Context, req department.DeleteDeptRequest) (interface{}, error) {
	return r.repository.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req department.DeleteDeptRequest) entity.DeleteDeptEntity {

	return entity.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}
