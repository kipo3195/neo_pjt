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

	var req chatService.SendChatRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 라인키 생성
	lineKey, sendDate := r.svc.LineKey.GetLineKey(ctx)

	// Message Broker에 publish
	input := adapter.MakeSendChatInput(sendUserHash, lineKey, sendDate, req.EventType, req.ChatSession, req.ChatRoom, req.ChatLine)
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
		LineKey:      lineKey,
		SendUserHash: input.ChatLine.SendUserHash,
		Cmd:          input.ChatLine.Cmd,
		Contents:     input.ChatLine.Contents,
		SendDate:     input.ChatLine.SendDate,
	}

	res := chatService.SendChatResponse{
		ChatRoom:    chatRoom,
		ChatLine:    chatLine,
		EventType:   input.EventType,
		ChatSession: input.ChatSession,
	}
	response.SendSuccess(c, res)
}
