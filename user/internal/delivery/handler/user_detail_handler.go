package handler

import (
	"user/internal/application/usecase"
)

type UserDetailHandler struct {
	usecase usecase.UserDetailUsecase
}

func NewUserDetailHandler(usecase usecase.UserDetailUsecase) *UserDetailHandler {
	return &UserDetailHandler{
		usecase: usecase,
	}
}
