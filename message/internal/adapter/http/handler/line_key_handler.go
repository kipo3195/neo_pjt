package handler

import (
	"message/internal/adapter/http/dto/lineKey"
	"message/internal/application/usecase"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
)

type LineKeyHandler struct {
	usecase usecase.LineKeyUsecase
}

func NewLineKeyHandler(usecase usecase.LineKeyUsecase) *LineKeyHandler {
	return &LineKeyHandler{
		usecase: usecase,
	}
}

func (h *LineKeyHandler) GetLineKey(c *gin.Context) {

	ctx := c.Request.Context()

	temp, _ := h.usecase.GetLineKey(ctx)

	res := lineKey.GetLineKeyResponse{
		LineKey: temp,
	}

	response.SendSuccess(c, res)

}
