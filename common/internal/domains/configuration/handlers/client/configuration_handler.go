package client

import (
	usecases "common/internal/domains/configuration/usecases/client"
)

type ConfigurationHandler struct {
	usecase usecases.ConfigurationUsecase
}

func NewConfigurationHandler(usecase usecases.ConfigurationUsecase) *ConfigurationHandler {
	return &ConfigurationHandler{
		usecase: usecase,
	}
}

// func (h *ConfigurationHandler) GetConfigHash(c *gin.Context) {

// 	// context 생성
// 	ctx := c.Request.Context()

// 	// 데이터 -> dto
// 	var req = requestDTO.GetConfigHashRequestBody{
// 		SkinHash:   c.Query("skinHash"),
// 		ConfigHash: c.Query("configHash"),
// 		Device:     c.Query("device"),
// 	}

// 	// 유효성 검증
// 	if req.SkinHash == "" || req.ConfigHash == "" || req.Device == "" {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
// 		return
// 	}

// 	// usecase 호출
// 	data := h.usecase.GetConfigHash(toConfigHashEntity(req), ctx)

// 	response.SendSuccess(c, data)

// }

// func toConfigHashEntity(requestDTO requestDTO.GetConfigHashRequestBody) entities.ConfigHashEntity {
// 	return entities.ConfigHashEntity{
// 		ConfigHash: requestDTO.ConfigHash,
// 		SkinHash:   requestDTO.SkinHash,
// 		Device:     requestDTO.Device,
// 	}
// }
