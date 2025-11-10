package handler

import (
	"message/internal/application/usecase"
	"message/internal/delivery/dto/lineKey"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
)

type LineKeyHandler struct {
	usecase usecase.LineKeyUsecase
}

func NewLineKeyHandler(uc usecase.LineKeyUsecase) *LineKeyHandler {
	return &LineKeyHandler{
		usecase: uc,
	}
}

func (h *LineKeyHandler) GetLineKey(c *gin.Context) {

	ctx := c.Request.Context()

	temp := h.usecase.GetLineKey(ctx)

	res := lineKey.GetLineKeyResponse{
		LineKey: temp,
	}

	response.SendSuccess(c, res)

}
