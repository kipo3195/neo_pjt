package appValidation

import (
	"core/internal/domain/appValidation/entity"
)

type AppValidationResponseBody struct {
	WorksCommonInfo *entity.WorksCommonInfo  `json:"worksCommonInfo"`
	WorksInfo       *entity.DeviceInitResult `json:"worksInfo"`
}

type AppValidationResponseDTO struct {
	Body AppValidationResponseBody
}
