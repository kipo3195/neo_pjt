package handler

import (
	"log"
	"message/internal/adapter/http/dto/lineKey"
	"message/internal/application/usecase"
	response "message/pkg/response"
	"time"

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
	log.Println("라인키 발급 전 sleep")
	time.Sleep(15 * time.Second)
	temp, _ := h.usecase.GetLineKey(ctx)
	log.Println("라인키 발급 temp")

	res := lineKey.GetLineKeyResponse{
		LineKey: temp,
	}

	response.SendSuccess(c, res)

}
