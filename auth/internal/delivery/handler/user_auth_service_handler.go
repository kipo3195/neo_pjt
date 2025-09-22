package handler

import (
	"auth/internal/application/orchestrator"

	"github.com/gin-gonic/gin"
)

type UserAuthServiceHandler struct {
	svc *orchestrator.UserAuthService
}

func NewUserAuthServiceHandler(svc *orchestrator.UserAuthService) *UserAuthServiceHandler {
	return &UserAuthServiceHandler{svc: svc}
}

func (h *UserAuthServiceHandler) UserAuthAndDeviceCheck(c *gin.Context) {

}
