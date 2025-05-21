package handlers

import (
	consts "auth/consts"
	"auth/dto"
	"auth/usecases"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	usecase usecases.AuthUsecase
}

func NewAuthHandler(uc usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) GetAuth(w http.ResponseWriter, r *http.Request) {

	// request 데이터 -> dto로 변경
	var authRequest dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 비즈니스 로직 호출
	Auth, err := h.usecase.GetAuth(authRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Entity -> dto로 변환은 handler에서 처리함.
	res := dto.AuthResponse{Result: Auth.Result, AccessToken: Auth.AccessToken, RefreshToken: Auth.RefreshToken, ConfigKey: Auth.ConfigKey}

	fmt.Println("Auth 결과 값 : ", res)
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) GenerateDeviceToken(w http.ResponseWriter, r *http.Request) {
	// response
	var res dto.GenerateDeviceTokenResponse

	// request의 header 데이터 -> dto로 변경
	header := &dto.GenerateDeviceTokenRequestHeader{
		Token: r.Header.Get("Authorization"),
	}

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
	var body dto.GenerateDeviceTokenRequest

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

	json.NewEncoder(w).Encode(res)

}
