package handlers

import (
	"core/config"
	consts "core/consts"
	"core/dto"
	"core/usecases"
	"encoding/json"
	"net/http"
)

type CoreHandler struct {
	usecase usecases.CoreUsecase
	sfg     *config.ServerConfig
}

func NewCoreHandler(sfg *config.ServerConfig, uc usecases.CoreUsecase) *CoreHandler {
	return &CoreHandler{usecase: uc, sfg: sfg}
}

// /app-validation
func (h *CoreHandler) GetAppValidation(w http.ResponseWriter, r *http.Request) {

	// 해당 API의 response
	var res = dto.AppValidationResponse{}

	// request의 header 데이터 -> dto로 변경
	var headerPrefix = h.sfg.ApiConfig.NeoHeaderPrefix
	header := &dto.AppValidationRequestHeader{
		Hash:   r.Header.Get(headerPrefix + "Hash"),
		Device: r.Header.Get(headerPrefix + "Device"),
		Uuid:   r.Header.Get(headerPrefix + "Uuid"),
	}

	if header.Hash == "" || header.Device == "" || header.Uuid == "" {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_104,
			Message: consts.E_104_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// request body 데이터 -> dto로 변경
	var body dto.AppValidationRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 배포 앱 hash 검증
	if h.usecase.CheckValidation(header) {

		// 클라이언트가 넘겨준 Domain : 테넌트 정보로 검증
		data, err := h.usecase.GetWorksInfo(body, header.Uuid)
		if err == nil {
			// http status code 200
			res.Code = consts.SUCCESS
			res.Data = data

		} else {
			// http status code 400
			res.Code = consts.ERROR
			res.Data = err
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// http status code 400
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.F_101,
			Message: consts.F_101_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(res)

}

// /config
func (h *CoreHandler) GetConfig(w http.ResponseWriter, r *http.Request) {

}
