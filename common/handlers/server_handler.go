package handlers

import (
	consts "common/consts"
	dto "common/dto/common"
	svDto "common/dto/server"
	"common/usecases"
	"context"
	"encoding/json"
	"fmt"
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
	var res = svDto.SvDeviceInitResponse{}

	// request의 header 데이터 -> dto로 변경
	header := &svDto.SvDeviceInitRequestHeader{
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
	var body *svDto.SvDeviceInitRequest
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
	data, err := h.usecase.DeviceInit(body, ctx)

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
