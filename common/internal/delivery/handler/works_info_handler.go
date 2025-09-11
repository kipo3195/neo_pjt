package handler

import "common/internal/application/usecase"

type WorksInfoHandler struct {
	usecase usecase.WorksInfoUsecase
}

func NewWorksInfoHandler(usecase usecase.WorksInfoUsecase) *WorksInfoHandler {
	return &WorksInfoHandler{
		usecase: usecase,
	}
}
