package handler

import "notificator/internal/application/usecase"

type LoginHandler struct {
	usecase usecase.LoginUsecase
}

func NewLoginHandler(usecase usecase.LoginUsecase) *LoginHandler {

	return &LoginHandler{
		usecase: usecase,
	}
}
