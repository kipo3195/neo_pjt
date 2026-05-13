package handler

import (
	"test/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	usecase usecase.ChatUsecase
}

func NewChatHandler(usecase usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
	}
}

func (r *ChatHandler) PutChatEvent(c *gin.Context) {

	// 채팅 발송
	// 더미 룸키 생성

}
