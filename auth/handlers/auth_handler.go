package handlers

import (
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
