package handler

import (
	"encoding/json"
	"message/internal/application/orchestrator"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/chatService"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatServiceHandler struct {
	svc *orchestrator.ChatService
}

func NewChatServiceHandler(svc *orchestrator.ChatService) *ChatServiceHandler {
	return &ChatServiceHandler{
		svc: svc,
	}
}

func (r *ChatServiceHandler) SendChat(c *gin.Context) {

	ctx := c.Request.Context()

	lineKey := r.svc.LineKey.GetLineKey(ctx)

	var req chatService.SendChatRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	input := adapter.MakeSendChatInput(lineKey, req.Contents, req.DestUsers)

	// 여기서 뽑아내는게 맞나?
	lineKey, err := r.svc.Chat.SendChat(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	res := chatService.SendChatResponse{
		LineKey: lineKey,
	}

	response.SendSuccess(c, res)

}
