package responseDTO

import (
	"core/internal/domains/appValidation/dto/external/commonResponseDTO"
	"core/internal/domains/appValidation/entities"
)

type AppValidationResponseBody struct {
	WorksCommonInfo *entities.WorksCommonInfo                 `json:"worksCommonInfo"`
	WorksInfo       *commonResponseDTO.DeviceInitResponseBody `json:"worksInfo"`
}

type AppValidationResponseDTO struct {
	Body AppValidationResponseBody
}
