package handler

import (
	"encoding/json"
	"message/internal/application/usecase"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/chatRoomTitle"
	"message/internal/delivery/util"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatRoomTitleHandler struct {
	usecase usecase.ChatRoomTitleUsecase
}

func NewChatRoomTitleHandler(usecase usecase.ChatRoomTitleUsecase) *ChatRoomTitleHandler {
	return &ChatRoomTitleHandler{
		usecase: usecase,
	}
}

func (h *ChatRoomTitleHandler) UpdateChatRoomTitle(c *gin.Context) {

	ctx := c.Request.Context()

	userHash := util.GetUserHashByAccessToken(c)

	var req chatRoomTitle.UpdateChatRoomTitleRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeUpdateChatRoomTitleInput(userHash, req.Org, req.RoomKey, req.Type, req.Title)
	output, err := h.usecase.UpdateChatRoomTitle(ctx, input)

	if err != nil {
		if err == consts.ErrDBresultNotFound {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F008, consts.MESSAGE_F008_MSG)
		} else if err == consts.ErrChatRoomKeyCheck {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F010, consts.MESSAGE_F010_MSG)
		} else if err == consts.ErrChatRoomTypeMismatch {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F011, consts.MESSAGE_F011_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	res := chatRoomTitle.UpdateChatRoomTitleResponse{
		UpdateDate: output,
	}

	response.SendSuccess(c, res)

}

func (h *ChatRoomTitleHandler) DeleteChatRoomTitle(c *gin.Context) {

	ctx := c.Request.Context()

	userHash := util.GetUserHashByAccessToken(c)

	var req chatRoomTitle.DeleteChatRoomTitleRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeDeleteChatRoomTitleInput(userHash, req.Org, req.RoomKey, req.Type)
	output, err := h.usecase.DeleteChatRoomTitle(ctx, input)

	if err != nil {
		if err == consts.ErrDBresultNotFound {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F009, consts.MESSAGE_F009_MSG)
		} else if err == consts.ErrChatRoomKeyCheck {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F010, consts.MESSAGE_F010_MSG)
		} else if err == consts.ErrChatRoomTypeMismatch {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F011, consts.MESSAGE_F011_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	res := chatRoomTitle.DeleteChatRoomTitleResponse{
		DeleteDate: output,
	}

	response.SendSuccess(c, res)

}
