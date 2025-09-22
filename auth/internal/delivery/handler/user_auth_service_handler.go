package handler

import "auth/internal/application/orchestrator"

type UserAuthServiceHandler struct {
	svc *orchestrator.UserAuthService
}

func NewUserAuthServiceHandler(svc *orchestrator.UserAuthService) *UserAuthServiceHandler {
	return &UserAuthServiceHandler{svc: svc}
}
