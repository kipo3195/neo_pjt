package repository

import (
	"core/internal/domain/appValidation/entity"
)

type AppValidationRepository interface {
	GetValidation(where entity.ValidationEntity) error
	GetWorksCommonInfo(worksCode string) (*entity.WorksCommonInfo, error)
}
