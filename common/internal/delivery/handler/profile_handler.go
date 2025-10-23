package handler

import "common/internal/application/usecase"

type ProfileHandler struct {
	usecase usecase.ProfileUsecase
}

func NewProfileHandler(usecase usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		usecase: usecase,
	}
}
