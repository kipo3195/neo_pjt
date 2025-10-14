package handler

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase/input"
	"auth/internal/consts"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"
	"errors"

	"auth/internal/delivery/dto/userAuthService"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserAuthServiceHandler struct {
	svc *orchestrator.UserAuthService
}

func NewUserAuthServiceHandler(svc *orchestrator.UserAuthService) *UserAuthServiceHandler {
	return &UserAuthServiceHandler{svc: svc}
}

func (h *UserAuthServiceHandler) UserAuthAndDeviceCheck(c *gin.Context) {

	ctx := c.Request.Context()
	var req userAuthService.UserAuthServiceRequest

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

	userAuthInput := input.MakeUserAuthInput(req.Id, req.Fv, req.Uuid)
	userAuthResult, err := h.svc.UserAuth.GetUserAuth(ctx, userAuthInput)

	if err != nil {
		if errors.Is(err, consts.ErrUserAuthFvMismatch) {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
		} else if errors.Is(err, consts.ErrUserAuthChallengeExpired) {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F007, consts.AUTH_F007_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	//var deviceRegistCheckOutput output.DeviceRegistStateOutput

	// 인증이 성공했을때, device 등록 체크
	if userAuthResult {

		deviceRegistCheckInput := input.MakeDeviceRegistCheckInput(req.Id, req.Uuid)
		challenge, err := h.svc.Device.GetDeviceRegistState(ctx, deviceRegistCheckInput)

		if err != nil {
			if err == consts.ErrDeviceNotRegist {
				// 등록되지 않은 device이므로 challenge 발급함.
				res := userAuthService.UserAuthServiceResponse{
					DeviceChallenge: challenge,
				}
				response.SendSuccess(c, res)
			} else {
				response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
			}
			return
		}

		// at, rt 발급 (token)
		userAuthTokenInput := input.MakeUserAuthTokenInput(req.Id, req.Uuid)
		h.svc.Token.GenerateAuthToken(ctx, userAuthTokenInput)

	} else {
		// 인증 실패.
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// res.AccessToken = deviceRegistCheckOutput.AccessToken
	// 	res.RefreshToken:    deviceRegistCheckOutput.RefreshToken,
	// 	res.DeviceChallenge: deviceRegistCheckOutput.DeviceRegistChallenge,
	// 	res.RefreshTokenExp: deviceRegistCheckOutput.RefreshTokenExp,
	// }

	response.SendSuccess(c, res)
}
