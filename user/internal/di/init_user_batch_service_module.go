package di

import (
	"user/internal/application/orchestrator"
	"user/internal/application/usecase"
	"user/internal/delivery/handler"
)

func InitUserBatchServiceModule(profile usecase.ProfileUsecase, userDetail usecase.UserDetailUsecase) *handler.UserBatchServiceHandler {
	service := orchestrator.NewUserBatchService(profile, userDetail)
	return handler.NewUserBatchServiceHandler(service)
}
