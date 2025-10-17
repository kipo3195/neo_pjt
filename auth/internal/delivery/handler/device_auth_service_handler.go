package handler

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase/input"
	"auth/internal/consts"
	"auth/internal/delivery/adapter"
	"auth/internal/delivery/dto/deviceAuthService"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"
	"log"

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

	deviceRegistInput := adapter.MakeDeviceRegistCheckInput(req.Id, req.Uuid, req.Challenge)
	deviceRegResult, err := h.svc.Device.DeviceRegistCheck(ctx, deviceRegistInput)

	if err != nil {
		// 등록에 따라 다르게 처리 필요 TODO
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

		// 발급 받았던 challenge 삭제
		removeDeviceChallengeInput := adapter.MakeRemoveDeviceChallengeInput(req.Id, req.Uuid)
		h.svc.Device.RemoveDeviceChallenge(ctx, removeDeviceChallengeInput)

		res := deviceAuthService.DeviceRegistResponse{
			RefreshToken:    output.RefreshToken,
			RefreshTokenExp: output.RefreshTokenExp,
			AccessToken:     output.AccessToken,
		}

		response.SendSuccess(c, res)
		return

	}

}

func (h *DeviceAuthServiceHandler) DeviceRefresh(c *gin.Context) {

	ctx := c.Request.Context()
	var req deviceAuthService.DeviceRefreshRequest

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

	// rt 만료 체크

	// device 등록 점검 (재발급 받을 수 있는지)
	deviceRegistInput := adapter.MakeDeviceRegistStateInput(req.Id, req.Uuid)
	challenge, err := h.svc.Device.GetDeviceRegistState(ctx, deviceRegistInput)

	log.Println("[DeviceRefresh] id : ", req.Id, ", new challenge : ", challenge)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}
	// device 정보 업데이트
	updateDeviceInfoInput := adapter.MakeUpdateDeviceInfoInput(req.Id, req.Uuid, req.ModelName, req.Version)
	err = h.svc.Device.UpdateDeviceInfo(ctx, updateDeviceInfoInput)

	if err != nil {
		if err == consts.ErrDeviceNotRegist {
			// 디바이스 등록되있지도 않은데 재발급 요청함
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F010, consts.AUTH_F010_MSG)
		} else {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	// at, rt 신규 생성
	generateAuthTokenInput := input.MakeGenerateAuthTokenInput(req.Id, req.Uuid)
	output, err := h.svc.Token.GenerateAuthToken(ctx, generateAuthTokenInput)
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	res := deviceAuthService.DeviceRefreshResponse{
		RefreshToken:    output.RefreshToken,
		RefreshTokenExp: output.RefreshTokenExp,
		AccessToken:     output.AccessToken,
	}

	// 발급 받았던 challenge 삭제
	removeDeviceChallengeInput := adapter.MakeRemoveDeviceChallengeInput(req.Id, req.Uuid)
	h.svc.Device.RemoveDeviceChallenge(ctx, removeDeviceChallengeInput)

	response.SendSuccess(c, res)

}
