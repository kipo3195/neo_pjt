package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
)

func InitAppValidationServiceModule(appValidation usecase.AppValidationUsecase, skin usecase.SkinUsecase, configuration usecase.ConfigurationUsecase) *handler.AppValidationServiceHandler {

	service := orchestrator.NewAppValidationService(appValidation, skin, configuration)

	return handler.NewAppValidationServiceHander(service)
}
