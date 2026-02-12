package handler

import (
	"context"
	"errors"
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
	log.Println("라인키 발급 전 슬립 15초 시작")
	time.Sleep(15 * time.Second)
	log.Println("라인키 발급 전 슬립 15초 끝")
	temp, _, err := h.usecase.GetLineKey(ctx)
	log.Println("라인키 : ", temp)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("deadLineExceed.")
			return
		}
	}

	res := lineKey.GetLineKeyResponse{
		LineKey: temp,
	}

	response.SendSuccess(c, res)

}
