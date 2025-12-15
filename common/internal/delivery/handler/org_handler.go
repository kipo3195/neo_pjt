package handler

import "common/internal/application/usecase"

type OrgHandler struct {
	usecase usecase.OrgUsecase
}

func NewOrgHandler(usecase usecase.OrgUsecase) *OrgHandler {

	return &OrgHandler{
		usecase: usecase,
	}
}
