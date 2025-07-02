package handlers

import (
	consts "auth/consts"
	dto "auth/dto/common"
	commonDto "auth/dto/server/common"
	"auth/usecases"
	"context"
	"encoding/json"
	"errors"
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

func (h *ServerHandler) AppTokenValidation(w http.ResponseWriter, r *http.Request) {
	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fmt.Println("1")
	// response
	var res dto.Response

	var req commonDto.AppTokenValidationRequest

	fmt.Println("2")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	data, err := h.usecase.AppTokenValidation(req, ctx)
	fmt.Println(data, err)
	// 이거 나중에 모듈화 꼭 할 것
	if err != nil || !data { // 에러
		fmt.Println(err)
		switch {
		case errors.Is(err, consts.ErrDbRowNotFound):
			// 매핑된 hash 정보가 없음
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.AUTH_F001, consts.AUTH_F001_MSG)
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, consts.ErrTokenExpired):
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_107, consts.E_107_MSG)
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, consts.ErrTokenSignatureInvalid):
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.AUTH_F005, consts.AUTH_F005_MSG)
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, consts.ErrDB):
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_102, consts.E_102_MSG)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, consts.ErrTokenParsing):
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_105, consts.E_105_MSG)
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, consts.ErrInvalidClaims):
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.AUTH_F002, consts.AUTH_F002_MSG)
			w.WriteHeader(http.StatusBadRequest)
		default:
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_500, consts.E_500_MSG)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		res.Result = consts.SUCCESS
		// 데이터가 있어야할까?
	}

	json.NewEncoder(w).Encode(res)
}
