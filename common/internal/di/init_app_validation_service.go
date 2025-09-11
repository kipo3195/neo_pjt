package di

import (
	"common/internal/application/orchestrator"
	"common/internal/application/usecase"
	"common/internal/delivery/handler"
)

func InitAppValidationHandler(appValidation usecase.AppValidationUsecase, skin usecase.SkinUsecase, configuration usecase.ConfigurationUsecase) *handler.AppValidationHandler {

	service := orchestrator.NewAppValidationService(appValidation, skin, configuration)

	return handler.NewAppValidationHander(service)
}
