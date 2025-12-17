package handler

import "admin/internal/application/usecase"

type UserAuthRegisterHandler struct {
	usecase usecase.UserAuthRegisterUsecase
}

func NewUserAuthRegisterHandler(usecase usecase.UserAuthRegisterUsecase) *UserAuthRegisterHandler {

	return &UserAuthRegisterHandler{
		usecase: usecase,
	}
}
