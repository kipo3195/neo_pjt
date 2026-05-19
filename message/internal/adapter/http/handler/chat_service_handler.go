package handler

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/adapter/http/dto/chatService"
	"message/internal/adapter/http/mapper"
	"message/internal/application/service"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatServiceHandler struct {
	svc *service.ChatService
}

func NewChatServiceHandler(svc *service.ChatService) *ChatServiceHandler {
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
	_, readDate, _ := r.svc.LineKey.GetLineKey(ctx)

	input := mapper.MakeReadChatInput(req.RoomKey, req.RoomType, userHash, readDate)
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
	lineKey, sendDate, _ := r.svc.LineKey.GetLineKey(ctx)

	if req.TestSender != "" {
		sendUserHash = req.TestSender
	}

	// Message Broker에 publish
	input := mapper.MakeSendChatInput(sendUserHash, lineKey, sendDate, req.EventType, req.ChatSession, req.ChatRoom, req.ChatLine, req.TransactionId)
	output, err := r.svc.Chat.SendChat(ctx, input)

	if err != nil {

		if err == consts.ErrPublishToMessageBrokerError {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F001, consts.MESSAGE_F001_MSG)
		} else if err == consts.ErrCacheResultNotFound {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F014, consts.MESSAGE_F014_MSG)
		} else if err == consts.ErrFileServiceGrpcCallErr {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F015, consts.MESSAGE_F015_MSG)
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
	for _, f := range output.SendChatFileOutput {

		temp := chatService.ChatFileData{
			FileId:       f.FileId,
			FileType:     f.FileType,
			FileName:     f.FileName,
			ThumbnailUrl: f.ThumbnailUrl,
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

func (r *ChatServiceHandler) SendBulkChat(c *gin.Context) {

	var req chatService.BulkSendChatRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	bulkInput := input.BulkSendChatInput{
		TotalCount:  req.TotalCount,
		PerSecond:   req.PerSecond,
		WorkerCount: req.WorkerCount,
		Contents:    req.Contents,
		Cmd:         req.Cmd,
		EventType:   req.EventType,
		ChatSession: req.ChatSession,
	}
	if bulkInput.WorkerCount == 0 {
		bulkInput.WorkerCount = 1
	}

	// 별도의 비동기 처리
	go func() {
		sentCount, err := r.svc.SendBulkChat(context.Background(), bulkInput)
		if err != nil {
			log.Println("[SendBulkChat] async failed. sent count :", sentCount, "err :", err)
			return
		}
		log.Println("[SendBulkChat] async finished. sent count :", sentCount)
	}()

	res := chatService.BulkSendChatResponse{
		TotalCount:  req.TotalCount,
		PerSecond:   req.PerSecond,
		WorkerCount: bulkInput.WorkerCount,
	}

	response.SendSuccess(c, res)
}
