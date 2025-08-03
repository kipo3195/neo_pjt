package handlers

import (
	"admin/usecases"
)

// 관리자 - org 서비스 연계 handler
type OrgHandler struct {
	usecase usecases.OrgUsecase
}

func NewOrgHandler(r usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: r}
}
