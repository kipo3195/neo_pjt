package appValidation

import (
	"core/internal/domain/appValidation/entity"
)

type AppValidationResponseBody struct {
	WorksCommonInfo *entity.WorksCommonInfo `json:"worksCommonInfo"`
	WorksInfo       *DeviceInitResponseBody `json:"worksInfo"`
}

type AppValidationResponseDTO struct {
	Body AppValidationResponseBody
}
