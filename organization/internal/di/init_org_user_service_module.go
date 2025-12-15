package di

import (
	"org/internal/application/orchestrator"
	"org/internal/application/usecase"
	"org/internal/delivery/handler"
)

func InitOrgUserServiceModule(org usecase.OrgUsecase, user usecase.UserUsecase) *handler.OrgUserServiceHandler {

	service := orchestrator.NewOrgUserService(org, user)
	return handler.NewOrgUserServiceHandler(service)
}
