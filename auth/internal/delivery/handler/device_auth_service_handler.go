package handler

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/adapter"
	"auth/internal/delivery/dto/deviceAuthService"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type DeviceAuthServiceHandler struct {
	svc *orchestrator.DeviceAuthService
}

func NewDeviceAuthServiceHandler(svc *orchestrator.DeviceAuthService) *DeviceAuthServiceHandler {
	return &DeviceAuthServiceHandler{svc: svc}
}

func (h *DeviceAuthServiceHandler) DeviceRegist(c *gin.Context) {

	ctx := c.Request.Context()
	var req deviceAuthService.DeviceRegistRequest

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

	deviceRegistInput := adapter.MakeDeviceRegistCheckInput(req.Id, req.Uuid)
	deviceRegResult, err := h.svc.Device.DeviceRegistCheck(ctx, deviceRegistInput)
	if err != nil {
		// 등록에 따라 다르게 처리 필요
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		// 이미 등록됨
		return
	}

	if deviceRegResult {
		// 등록 성공
		generateAuthTokenInput := input.MakeGenerateAuthTokenInput(req.Id, req.Uuid)
		output, err := h.svc.Token.GenerateAuthToken(ctx, generateAuthTokenInput)
		if err != nil {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
			return
		}

		res := deviceAuthService.DeviceRegistResponse{
			RefreshToken:    output.RefreshToken,
			RefreshTokenExp: output.RefreshTokenExp,
			AccessToken:     output.AccessToken,
		}

		response.SendSuccess(c, res)
		return

	}

}
