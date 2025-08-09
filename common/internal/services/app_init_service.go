package services

import (
	appValidationUsecase "common/internal/domains/appValidation/usecases/server"
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	skinUsecase "common/internal/domains/skin/usecases/client"
)

type AppInitService struct {
	validator     appValidationUsecase.AppValidationUsecase
	skin          skinUsecase.SkinUsecase
	configuration configurationUsecase.ConfigurationUsecase
}

func NewAppInitService(v appValidationUsecase.AppValidationUsecase, s skinUsecase.SkinUsecase, c configurationUsecase.ConfigurationUsecase) *AppInitService {
	return &AppInitService{v, s, c}
}
