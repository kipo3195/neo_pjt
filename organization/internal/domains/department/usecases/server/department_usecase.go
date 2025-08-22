package server

import (
	repositories "org/internal/domains/department/repositories/server"
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
