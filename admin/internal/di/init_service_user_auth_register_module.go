package di

import (
	"admin/internal/application/orchestrator"
	"admin/internal/application/usecase"
	"admin/internal/delivery/handler"
)

func InitServiceUserAuthRegisterServiceModule(serviceUser usecase.ServiceUserUsecase, userAuthRegister usecase.UserAuthRegisterUsecase) *handler.ServiceUserAuthRegisterHandler {
	service := orchestrator.NewServiceUserAuthRegisterService(serviceUser, userAuthRegister)
	return handler.NewServiceUserAuthRegisterHandler(service)

}
