package handlers

import (
	"core/dto"
	"core/usecases"
	"encoding/json"
	"net/http"
)

type CoreHandler struct {
	usecase usecases.CoreUsecase
}

func NewCoreHandler(uc usecases.CoreUsecase) *CoreHandler {
	return &CoreHandler{usecase: uc}
}

func (h *CoreHandler) GetValidation(w http.ResponseWriter, r *http.Request) {

	header, body, err := h.usecase.GetValidationData(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var res = dto.ValidationResponse{}
	if h.usecase.CheckValidation(header) {
		data, err := h.usecase.GetWorksInfo(body)
		if err != nil {
			res.Code = 200
			res.Data = "일치하는 도메인 없음. "
		} else {
			res.Code = 200
			res.Data = data

		}
		json.NewEncoder(w).Encode(res)
	} else {
		res.Code = 200
		res.Data = "앱 검증 실패"
	}
}
