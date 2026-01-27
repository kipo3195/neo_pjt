package handler

import (
	"encoding/json"
	"message/internal/application/orchestrator"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/chatLineService"
	"message/internal/delivery/util"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatLineServiceHandler struct {
	svc *orchestrator.ChatLineService
}

func NewChatLineServiceHandler(svc *orchestrator.ChatLineService) *ChatLineServiceHandler {
	return &ChatLineServiceHandler{
		svc: svc,
	}
}

func (r *ChatLineServiceHandler) GetChatLineEvent(c *gin.Context) {
	ctx := c.Request.Context()

	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatLineService.GetChatLineEventRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeGetChatLineEventInput(reqUserHash, req.Org, req.RoomType, req.RoomKey, req.LineKey)
	out, err := r.svc.Chat.GetChatLineEvent(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	data := make([]chatLineService.ChatLIneEventDto, 0)

	for _, o := range out {

		temp := chatLineService.ChatLIneEventDto{
			EventType:     o.EventType,
			Cmd:           o.Cmd,
			LineKey:       o.LineKey,
			TargetLineKey: o.TargetLineKey,
			Contents:      o.Contents,
			SendUserHash:  o.SendUserHash,
			SendDate:      o.SendDate,
			FileId:        o.FileId,
			FileName:      o.FileName,
			FileType:      o.FileType,
			ThumbnailUrl:  o.ThumbnailUrl,
		}
		data = append(data, temp)
	}

	res := chatLineService.GetChatLineEventResponse{
		ChatLineEventData: data,
	}

	response.SendSuccess(c, res)
}
