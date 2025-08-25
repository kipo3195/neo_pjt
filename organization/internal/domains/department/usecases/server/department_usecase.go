package server

import (
	"context"
	"log"
	"org/entities"
	requestDTO "org/internal/domains/department/dto/server/requestDTO"
	repositories "org/internal/domains/department/repositories/server"
	"org/internal/utils"
)

type departmentUsecase struct {
	repository repositories.DepartmentRepository
}

type DepartmentUsecase interface {
	CreateDept(ctx context.Context, requestDTO requestDTO.CreateDeptRequest) (interface{}, error)

	DeleteDept(ctx context.Context, req requestDTO.DeleteDeptRequest) (interface{}, error)

	CreateDeptUser(ctx context.Context, req requestDTO.CreateDeptUserRequest) (interface{}, error)
	DeleteDeptUser(ctx context.Context, req requestDTO.DeleteDeptUserRequest) (interface{}, error)
}

func NewDepartmentUsecase(repository repositories.DepartmentRepository) DepartmentUsecase {
	return &departmentUsecase{
		repository: repository,
	}
}

func (r *departmentUsecase) CreateDept(ctx context.Context, req requestDTO.CreateDeptRequest) (interface{}, error) {
	return r.repository.PutDept(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(requestDTO requestDTO.CreateDeptRequest) entities.CreateDeptEntity {

	return entities.CreateDeptEntity{
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

func (r *departmentUsecase) DeleteDeptUser(ctx context.Context, req requestDTO.DeleteDeptUserRequest) (interface{}, error) {
	return r.repository.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req requestDTO.DeleteDeptUserRequest) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *departmentUsecase) CreateDeptUser(ctx context.Context, req requestDTO.CreateDeptUserRequest) (interface{}, error) {

	updateHash := utils.MakeUpdateHash()
	log.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repository.PutDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req requestDTO.CreateDeptUserRequest, updateHash string) entities.CreateDeptUserEntity {

	return entities.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		DeptOrg:              req.DeptOrg,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
		UpdateHash:           updateHash,
	}
}

func (r *departmentUsecase) DeleteDept(ctx context.Context, req requestDTO.DeleteDeptRequest) (interface{}, error) {
	return r.repository.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req requestDTO.DeleteDeptRequest) entities.DeleteDeptEntity {

	return entities.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}
