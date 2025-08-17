package server

import (
	worksInfoUsecase "common/internal/domains/worksInfo/usecases/server"
)

type WorksInfoHandler struct {
	usecase worksInfoUsecase.WorksInfoUsecase
}

func NewWorksInfoHandler(usecase worksInfoUsecase.WorksInfoUsecase) *WorksInfoHandler {
	return &WorksInfoHandler{
		usecase: usecase,
	}
}
