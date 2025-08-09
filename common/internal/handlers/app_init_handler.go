package handlers

import (
	"net/http"

	"common/internal/services"

	"github.com/gin-gonic/gin"
)

type AppInitHandler struct {
	svc *services.AppInitService
}

func NewAppInitHander(svc *services.AppInitService) *AppInitHandler {
	return &AppInitHandler{svc: svc}
}

// POST /server/v1/app-validation
func (h *AppInitHandler) GetAppValidation(c *gin.Context) {

	// 실제 비즈니스 로직 처리? svc를 호출 기존 handler와 동일하게 처리하도록 수정 필요.
	var req struct {
		AppID      string `json:"app_id" binding:"required"`
		DeviceInfo string `json:"device_info" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.svc.GetAppProfile(req.AppID, req.DeviceInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
