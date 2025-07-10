package response

import (
	commonSvResDto "core/dto/server/common/response"
	"core/entities"
)

type AppValidationResponseBody struct {
	WorksCommonInfo *entities.WorksCommonInfo              `json:"worksCommonInfo"`
	WorksInfo       *commonSvResDto.DeviceInitResponseBody `json:"worksInfo"`
}

type AppValidationResponseDTO struct {
	Body AppValidationResponseBody
}
