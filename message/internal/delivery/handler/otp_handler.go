package handler

import (
	"encoding/json"
	"log"
	"message/internal/application/usecase"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/otp"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	usecase usecase.OtpUsecase
}

func NewOtpHandler(uc usecase.OtpUsecase) *OtpHandler {
	return &OtpHandler{
		usecase: uc,
	}
}

func (h *OtpHandler) OtpKeyRegist(c *gin.Context) {

	ctx := c.Request.Context()

	var req otp.OtpKeyRegistRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	input := adapter.MakeOtpKeyRegistInput(req.Id, req.Uuid, req.ChKey, req.NoKey)
	output, err := h.usecase.OtpKeyRegist(ctx, input)

	if err != nil {
		log.Println("OtpKeyRegist err : ", err)
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	res := otp.OtpKeyRegistResponse{
		OtpRegDate:   output.OtpRegDate,
		SvKeyVersion: output.SvKeyVersion,
	}

	response.SendSuccess(c, res)
}
