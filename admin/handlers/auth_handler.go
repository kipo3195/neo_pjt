package handlers

import "admin/usecases"

type AuthHandler struct {
	usecase usecases.AuthUsecase
}

func NewAuthHandler(r usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: r}
}
