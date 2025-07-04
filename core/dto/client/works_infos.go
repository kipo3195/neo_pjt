package dto

import (
	commonDto "core/dto/server/common"
	"core/entities"
)

type WorksInfos struct {
	WorksCommonInfo *entities.WorksCommonInfo     `json:"worksCommonInfo"`
	WorksInfo       *commonDto.DeviceInitResponse `json:"worksInfo"`
}
