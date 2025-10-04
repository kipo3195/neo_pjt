package orchestrator

import "org/internal/application/usecase"

type DummyDataService struct {
	Department usecase.DepartmentUsecase
	Org        usecase.OrgUsecase
	User       usecase.UserUsecase
}

func NewDummyDataService(u usecase.UserUsecase, d usecase.DepartmentUsecase, o usecase.OrgUsecase) *DummyDataService {
	return &DummyDataService{
		Department: d,
		Org:        o,
		User:       u,
	}
}
