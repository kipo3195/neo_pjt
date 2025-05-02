package handlers

import (
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
