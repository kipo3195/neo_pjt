package di

import (
	"org/internal/application/orchestrator"
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
)

func InitDummyDataServiceModule(department usecase.DepartmentUsecase, org usecase.OrgUsecase, user usecase.UserUsecase) *handler.DummyDataServiceHandler {

	service := orchestrator.NewDummyDataService(user, department, org)
	return handler.NewDummyDataServiceHandler(service)
}
