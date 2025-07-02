package handlers

import (
	clDto "common/dto/client"
	dto "common/dto/common"
	"common/usecases"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	consts "common/consts"

	"github.com/go-playground/validator"
)

type PublicHandler struct {
	usecase usecases.PublicUsecase
}

func NewPublicHandler(uc usecases.PublicUsecase) *PublicHandler {
	return &PublicHandler{usecase: uc}
}

func (h *PublicHandler) AppValidation(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response
	var res dto.Response

	fmt.Println("1")
	// request body 데이터 -> dto로 변경
	var req = clDto.AppValidationRequest{
		Uuid:       r.URL.Query().Get("uuid"),
		AppToken:   r.URL.Query().Get("appToken"),
		Device:     r.URL.Query().Get("device"),
		SkinHash:   r.URL.Query().Get("skinHash"),
		ConfigHash: r.URL.Query().Get("configHash"),
	}

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	res.Result = consts.ERROR
	// 	res.Data = dto.ErrorResponse{
	// 		Code:    consts.E_103,
	// 		Message: consts.E_103_MSG,
	// 	}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(res)
	// 	return
	// }
	fmt.Println("2")
	// 유효성 검증 로직
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		w.WriteHeader(http.StatusBadRequest)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_108,
			Message: consts.E_108_MSG,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Println("3")
	// 검증
	data, err := h.usecase.AppValidation(ctx, req)

	if err != nil || !data {
		switch {
		case errors.Is(err, consts.ErrInvalidClaims):
			w.WriteHeader(http.StatusUnauthorized)
			res.Result = consts.ERROR
			res.Data = dto.ErrorResponse{
				Code:    consts.E_106,
				Message: consts.E_106_MSG,
			}
		case errors.Is(err, consts.ErrSkinHashInvalid):
			w.WriteHeader(http.StatusBadRequest)
			res.Result = consts.FAIL
			res.Data = dto.ErrorResponse{
				Code:    consts.COMMON_F001,
				Message: consts.COMMON_F001_MSG,
			}
		case errors.Is(err, consts.ErrConfigHashInvalid):
			w.WriteHeader(http.StatusBadRequest)
			res.Result = consts.FAIL
			res.Data = dto.ErrorResponse{
				Code:    consts.COMMON_F002,
				Message: consts.COMMON_F002_MSG,
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			res.Result = consts.ERROR
			res.Data = dto.ErrorResponse{
				Code:    consts.E_500,
				Message: consts.E_500_MSG,
			}
		}

	} else if data {
		res.Result = consts.SUCCESS
		res.Data = dto.ErrorResponse{}
	}

	json.NewEncoder(w).Encode(res)
}
