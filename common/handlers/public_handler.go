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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	defer r.Body.Close()

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
			writeJSON(w, http.StatusOK, res)
		case errors.Is(err, consts.ErrSkinHashInvalid):
			w.WriteHeader(http.StatusBadRequest)
			res.Result = consts.FAIL
			res.Data = dto.ErrorResponse{
				Code:    consts.COMMON_F001,
				Message: consts.COMMON_F001_MSG,
			}
			writeJSON(w, http.StatusOK, res)
		case errors.Is(err, consts.ErrConfigHashInvalid):
			w.WriteHeader(http.StatusBadRequest)
			res.Result = consts.FAIL
			res.Data = dto.ErrorResponse{
				Code:    consts.COMMON_F002,
				Message: consts.COMMON_F002_MSG,
			}
			writeJSON(w, http.StatusOK, res)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			res.Result = consts.ERROR
			res.Data = dto.ErrorResponse{
				Code:    consts.E_500,
				Message: consts.E_500_MSG,
			}
			writeJSON(w, http.StatusOK, res)
		}
		return
	} else if data {
		res.Result = consts.SUCCESS
		res.Data = dto.ErrorResponse{}
	}

	json.NewEncoder(w).Encode(res)
}

func (h *PublicHandler) AppTokenRefresh(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	defer r.Body.Close()

	var res clDto.AppTokenRefreshResponse
	var body clDto.AppTokenRefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.Result = consts.FAIL
		res.Data = newErrorResp(consts.E_103, consts.E_103_MSG)
		writeJSON(w, http.StatusBadRequest, res)
		return
	}

	data, err := h.usecase.AppTokenReIssue(ctx, body)

	if err != nil {
		switch {
		case errors.Is(err, consts.ErrRefreshTokenAuthInvalid):
			// 토큰 검증 실패
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.COMMON_F003, consts.COMMON_F003_MSG)
			writeJSON(w, http.StatusBadRequest, res)
		case errors.Is(err, consts.ErrRefreshTokenAuthExpired):
			// 토큰 만료
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.COMMON_F004, consts.COMMON_F004_MSG)
			writeJSON(w, http.StatusBadRequest, res)
		default:
			// 서버 에러
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_500, consts.E_500_MSG)
			writeJSON(w, http.StatusInternalServerError, res)
		}
		return
	}

	res.Result = consts.SUCCESS
	res.Data = data
	writeJSON(w, http.StatusOK, res)
}

func newErrorResp(code, msg string) *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Code:    code,
		Message: msg,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// 로깅 권장
		fmt.Printf("json encoding failed: %v\n", err)
	}
}
