package handler

import (
	"encoding/json"
	"message/internal/application/orchestrator"
	"message/internal/consts"
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

	// 사용자 정보 파싱
	hash := c.Value(consts.USER_HASH)
	sendUserHash, ok := hash.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	// 라인키 조회
	lineKey := r.svc.LineKey.GetLineKey(ctx)
	var req chatService.SendChatRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// Message Broker에 publish
	input := adapter.MakeSendChatInput(sendUserHash, req.EventType, lineKey, req.ChatRoom, req.ChatLine)
	// err := r.svc.Chat.SendChat(ctx, input)

	// if err != nil {
	// 	if err == consts.ErrPublishToMessageBrokerError {
	// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F001, consts.MESSAGE_F001_MSG)
	// 	} else {
	// 		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	// 	}
	// 	return
	// }

	chatRoom := chatService.ChatRoomData{
		RoomKey:  input.ChatRoom.RoomKey,
		RoomType: input.ChatRoom.RoomType,
	}

	chatLine := chatService.ChatLineData{
		SendUserHash: sendUserHash,
		LineKey:      lineKey,
		Contents:     input.ChatLine.Contents,
		EventType:    input.ChatLine.EventType,
	}

	res := chatService.SendChatResponse{
		ChatRoom: chatRoom,
		ChatLine: chatLine,
	}
	response.SendSuccess(c, res)
}
