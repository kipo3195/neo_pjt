package di

import (
	"user/internal/application/orchestrator"
	"user/internal/application/usecase"
	"user/internal/delivery/handler"
)

func InitUserInfoServiceModule(profile usecase.ProfileUsecase, userDetail usecase.UserDetailUsecase) *handler.UserInfoServiceHandler {
	service := orchestrator.NewUserInfoService(profile, userDetail)
	return handler.NewUserInfoServiceHandler(service)
}
