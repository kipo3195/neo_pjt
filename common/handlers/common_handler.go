package handlers

import (
	consts "common/consts"
	clDto "common/dto/client"
	dto "common/dto/common"
	"common/entities"
	"common/usecases"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) GetConfigHash(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// 데이터 -> dto
	var req = clDto.GetConfigHash{
		SkinHash:   c.Query("skinHash"),
		ConfigHash: c.Query("configHash"),
		Device:     c.Query("device"),
	}

	// 유효성 검증
	if req.SkinHash == "" || req.ConfigHash == "" || req.Device == "" {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data := h.usecase.GetConfigHash(toConfigHashEntity(req), ctx)

	res.Result = consts.SUCCESS
	res.Data = data

	// response.
	json.NewEncoder(w).Encode(res)

}

func toConfigHashEntity(dto clDto.GetConfigHash) entities.ConfigHashEntity {
	return entities.ConfigHashEntity{
		ConfigHash: dto.ConfigHash,
		SkinHash:   dto.SkinHash,
		Device:     dto.Device,
	}
}
