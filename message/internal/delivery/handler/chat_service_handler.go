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
func (r *ChatServiceHandler) ReadChat(c *gin.Context) {
	ctx := c.Request.Context()

	// 사용자 정보 파싱
	hash := c.Value(consts.USER_HASH)
	userHash, ok := hash.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatService.ReadChatRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 라인키 생성 -> 사실상 response의 sendDate 값을 구하기 용도
	_, readDate := r.svc.LineKey.GetLineKey(ctx)

	input := adapter.MakeReadChatInput(req.RoomKey, req.RoomType, userHash, readDate)
	err := r.svc.Chat.ReadChat(ctx, input)

	if err != nil {
		if err == consts.ErrPublishToMessageBrokerError {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F001, consts.MESSAGE_F001_MSG)
		} else if err == consts.ErrDBResultNotUpdate {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F013, consts.MESSAGE_F013_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	res := chatService.ReadChatResponse{
		ReadDate: readDate,
	}

	response.SendSuccess(c, res)

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

	if req.TestSender != "" {
		sendUserHash = req.TestSender
	}

	// Message Broker에 publish
	input := adapter.MakeSendChatInput(sendUserHash, lineKey, sendDate, req.EventType, req.ChatSession, req.ChatRoom, req.ChatLine, req.ChatFile)
	output, err := r.svc.Chat.SendChat(ctx, input)

	if err != nil {
		if err == consts.ErrPublishToMessageBrokerError {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F001, consts.MESSAGE_F001_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	chatRoom := chatService.ChatRoomData{
		RoomKey:    input.ChatRoom.RoomKey,
		RoomType:   input.ChatRoom.RoomType,
		SecretFlag: input.ChatRoom.SecretFlag,
	}

	chatLine := chatService.ChatLineData{
		LineKey:       input.ChatLine.LineKey,
		SendUserHash:  input.ChatLine.SendUserHash,
		Cmd:           input.ChatLine.Cmd,
		Contents:      input.ChatLine.Contents,
		SendDate:      input.ChatLine.SendDate,
		TargetLineKey: input.ChatLine.TargetLineKey,
	}

	chatFile := make([]chatService.ChatFileData, 0)
	for _, f := range input.ChatFile {

		var thumbnailUrl string

		// 생성된 썸네일 url이 있는 경우
		v, exists := output.ThumbnailMap[f.FileId]
		if exists {
			thumbnailUrl = v
		}

		temp := chatService.ChatFileData{
			FileId:       f.FileId,
			FileExt:      f.FileExt,
			FileName:     f.FileName,
			ThumbnailUrl: thumbnailUrl,
		}
		chatFile = append(chatFile, temp)
	}

	res := chatService.SendChatResponse{
		ChatRoom:    chatRoom,
		ChatLine:    chatLine,
		EventType:   input.EventType,
		ChatSession: input.ChatSession,
		ChatFile:    chatFile,
	}

	response.SendSuccess(c, res)
}
