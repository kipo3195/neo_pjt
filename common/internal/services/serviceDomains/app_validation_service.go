package servicesDomains

import (
	appValidationUsecase "common/internal/domains/appValidation/usecases/server"
	configurationUsecase "common/internal/domains/configuration/usecases/client"
	skinUsecase "common/internal/domains/skin/usecases/client"
)

type AppValidationService struct {
	Validator     appValidationUsecase.AppValidationUsecase
	Skin          skinUsecase.SkinUsecase
	Configuration configurationUsecase.ConfigurationUsecase
}

func NewAppValidationService(v appValidationUsecase.AppValidationUsecase, s skinUsecase.SkinUsecase, c configurationUsecase.ConfigurationUsecase) *AppValidationService {
	return &AppValidationService{v, s, c}
}
