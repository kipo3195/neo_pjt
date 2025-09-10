package orchestrator

import (
	"common/internal/application/usecase"
)

type AppValidationService struct {
	Validator     usecase.AppValidationUsecase
	Skin          usecase.SkinUsecase
	Configuration usecase.ConfigurationUsecase
}

func NewAppValidationService(v usecase.AppValidationUsecase, s usecase.SkinUsecase, c usecase.ConfigurationUsecase) *AppValidationService {
	return &AppValidationService{v, s, c}
}
