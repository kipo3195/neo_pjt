package handlers

import (
	consts "common/consts"
	dto "common/dto/common"
	adminDto "common/dto/server/admin"
	commonDto "common/dto/server/common"
	"common/usecases"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServerHandler struct {
	usecase usecases.ServerUsecase
}

func NewServerHandler(uc usecases.ServerUsecase) *ServerHandler {
	return &ServerHandler{usecase: uc}
}

func (h *ServerHandler) DeviceInit(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 해당 API의 response
	var res = commonDto.DeviceInitResponse{}

	// core 서비스에서 온 토큰 검증은 미들 웨어에서

	// request body 데이터 -> dto로 변경
	var body *commonDto.DeviceInitRequest
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

	fmt.Printf("CORE 서버에서 호출, uuid : %s, worksCode : %s \n", body.Uuid, body.WorksCode)

	// DB에서 해당 works의 정보조회 + AUTH에서 토큰 발급 요청
	data, err := h.usecase.DeviceInit(ctx, body)

	if err != nil {
		res.Code = consts.FAIL
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res.Code = consts.SUCCESS
		res.Data = data
	}

	fmt.Printf("res.Data 타입: %T\n", res.Data)
	jsonBytes, _ := json.Marshal(res)
	log.Println("최종 JSON 응답:", string(jsonBytes))

	json.NewEncoder(w).Encode(res)

}

