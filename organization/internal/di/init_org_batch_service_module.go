package di

import (
	"org/internal/application/orchestrator"
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
)

func InitOrgBatchServiceModule(department usecase.DepartmentUsecase, org usecase.OrgUsecase, user usecase.UserUsecase) *handler.OrgBatchServiceHandler {

	service := orchestrator.NewOrgBatchService(department, org, user)

	return handler.NewOrgBatchServiceHandler(service)
}
