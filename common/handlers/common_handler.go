package handlers

import (
	consts "common/consts"
	clDto "common/dto/client"
	dto "common/dto/common"
	"common/entities"
	"common/usecases"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) GetConfigHash(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = dto.Response{}

	// 데이터 -> dto
	var req = clDto.GetConfigHash{
		SkinHash:   r.URL.Query().Get("skinHash"),
		ConfigHash: r.URL.Query().Get("configHash"),
		Device:     r.URL.Query().Get("device"),
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
