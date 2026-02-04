package handler

import (
	"encoding/json"
	"log"
	"message/internal/adapter/http/dto/otp"
	"message/internal/adapter/http/mapper"
	"message/internal/adapter/util"
	"message/internal/application/usecase"
	"message/internal/consts"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// id 기반으로 등록하는 이유는 userHash가 너무 길기 때문이다 + device 등록 시점에 otp를 등록하는데 해당 시점에는 AT가 없으므로 hash를 뽑아낼 수 없다.
	input := mapper.MakeOtpKeyRegistInput(req.Id, req.Uuid, req.DevicePubKey)
	output, err := h.usecase.OtpKeyRegist(ctx, input)

	if err != nil {
		log.Println("OtpKeyRegist err : ", err)
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	res := otp.OtpKeyRegistResponse{
		OtpRegDate:       output.OtpRegDate,
		SvChatKeyVersion: output.SvChatKeyVersion,
		SvNoteKeyVersion: output.SvNoteKeyVersion,
	}

	response.SendSuccess(c, res)
}

func (h *OtpHandler) GetMyOtpInfo(c *gin.Context) {

	ctx := c.Request.Context()

	// 등록도 ID 기반으로 하기 때문.
	userId := util.GetUserHashByAccessToken(c)

	var req otp.MyOtpInfoRequest

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

	input := mapper.MakeMyOtpInfoInput(userId, req.Uuid, req.VersionType, req.VersionInfo)
	output, err := h.usecase.GetMyOtpInfo(ctx, input)

	if err != nil {
		log.Println("GetMyOtpInfo err : ", err)
		if err == consts.ErrDBresultNotFound {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F002, consts.MESSAGE_F002_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	res := otp.MyOtpInfoResponse{
		MyOtpInfo: make([]otp.MyOtpInfo, 0), // ← 값 없어도 빈 배열 [] 보장
	}

	for _, info := range output {
		myOtpInfo := otp.MyOtpInfo{
			Version:    info.Version,
			KeyType:    info.KeyType,
			Key:        info.Key,
			OtpRegDate: info.OtpRegDate,
		}
		res.MyOtpInfo = append(res.MyOtpInfo, myOtpInfo)
	}
	response.SendSuccess(c, res)
}
