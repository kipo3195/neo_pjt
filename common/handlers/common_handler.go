package handlers

import (
	consts "common/consts"
	"common/dto"
	"common/usecases"
	"encoding/json"
	"fmt"
	"net/http"
)

type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) DeviceInit(w http.ResponseWriter, r *http.Request) {

	// 해당 API의 response
	var res = dto.DeviceInitResponse{}

	// request의 header 데이터 -> dto로 변경
	header := &dto.DeviceInitRequestHeader{
		Token: r.Header.Get("Authorization"), // const로 TODO
	}

	fmt.Println("CORE 서버에서 호출, 토큰 정보 : ", header.Token)

	if header.Token == "" {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_104,
			Message: consts.E_104_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// core서비스에서 온 토큰 검증 필요 todo

	// request body 데이터 -> dto로 변경
	var body *dto.DeviceInitRequest
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

	fmt.Printf("CORE 서버에서 호출, uuid : %s, domain : %s \n", body.Uuid, body.Domain)

	// DB에서 해당 works의 정보조회 + AUTH에서 토큰 발급 요청
	data, err := h.usecase.DeviceInit(body)

	if err != nil {
		res.Code = consts.FAIL
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res.Code = consts.SUCCESS
		res.Data = data
	}

	json.NewEncoder(w).Encode(res)

}

func (h *CommonHandler) GetConfig(w http.ResponseWriter, r *http.Request) {

	// request 데이터 -> dto로 변경
	var configRequest dto.ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&configRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("클라이언트 요청 수신 : ", configRequest)

	// 비즈니스 로직 호출
	config, err := h.usecase.GetConfig(configRequest)

	if err != nil {
		fmt.Println("없는 파일 요청.. ")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 파일 다운로드 응답
	w.Header().Set("Content-Disposition", "attachment; filename=\""+config.FileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(config.Content)
}
