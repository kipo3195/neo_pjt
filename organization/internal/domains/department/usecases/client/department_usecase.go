package client

import (
	repositories "org/internal/domains/department/repositories/client"
)

type departmentUsecase struct {
	repository repositories.DepartmentRepository
}

type DepartmentUsecase interface {
}

func NewDepartmentUsecase(repository repositories.DepartmentRepository) DepartmentUsecase {
	return &departmentUsecase{
		repository: repository,
	}
}
