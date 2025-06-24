package handlers

import (
	consts "auth/consts"
	clDto "auth/dto/client"
	dto "auth/dto/common"
	svDto "auth/dto/server"
	"auth/usecases"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	usecase usecases.AuthUsecase
}

func NewAuthHandler(uc usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func newErrorResp(code, msg string) *dto.ErrorResponse {
	return &dto.ErrorResponse{
		Code:    code,
		Message: msg,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// response
	var res dto.Response

	// request의 header 데이터 -> dto로 변경
	header := &clDto.LoginRequestHeader{
		Token: r.Header.Get("X-NEO-AuthToken"),
		Uuid:  r.Header.Get("X-NEO-Uuid"),
	}
	// header 검증
	if header.Token == "" {
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_104,
			Message: consts.E_104_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// request body 데이터 -> dto로 변경
	var body *clDto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 비즈니스 로직 호출
	Auth, err, failFlag := h.usecase.GetAuth(header, body)

	if failFlag { // 인증 실패
		res.Result = consts.FAIL
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	} else if err != nil { // 에러
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Entity -> dto로 변환은 handler에서 처리함.
		res.Result = consts.SUCCESS
		res.Data = clDto.AuthResponse{
			AccessToken:  Auth.AccessToken,
			RefreshToken: Auth.RefreshToken,
		}
	}
	json.NewEncoder(w).Encode(res)

}

func (h *AuthHandler) GenerateDeviceToken(w http.ResponseWriter, r *http.Request) {
	// response
	var res svDto.SvGenerateDeviceTokenResponse

	// request의 header 데이터 -> dto로 변경
	header := &svDto.SvGenerateDeviceTokenRequestHeader{
		Token: r.Header.Get("Authorization"),
	}

	fmt.Println("common service에서 호출시 던진 토큰 ", header.Token)

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
	// 서버의 토큰 검증 필요

	// request body 데이터 -> dto로 변경
	var body svDto.SvGenerateDeviceTokenRequest

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

	// 토큰 발급, DB 저장.
	result, err := h.usecase.GenerateDeviceToken(body)
	fmt.Println("handler에서 토큰 구조체 반환 result : ", result)

	if err != nil {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		res.Code = consts.SUCCESS
		res.Data = result
	}

	fmt.Println("handler에서 결과 반환 res : ", res)

	json.NewEncoder(w).Encode(res)

}

func (h *AuthHandler) AppTokenValidation(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response
	var res dto.Response

	fmt.Println("1")
	// request body 데이터 -> dto로 변경
	var req = &clDto.AppTokenValidationRequest{
		Uuid:     r.URL.Query().Get("uuid"),
		AppToken: r.URL.Query().Get("appToken"),
	}

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
	fmt.Println("2")
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
	data, err := h.usecase.AppTokenValidation(ctx, req)

	if err != nil || !data { // 에러
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
