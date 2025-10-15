package handler

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase/input"
	"auth/internal/consts"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"
	"errors"
	"log"

	"auth/internal/delivery/adapter"
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

	// 인증 도메인
	userAuthInput := input.MakeUserAuthInput(req.Id, req.Fv, req.Uuid)
	userAuthResult, err := h.svc.UserAuth.GetUserAuth(ctx, userAuthInput)

	if err != nil {
		log.Println("[UserAuthAndDeviceCheck] 인증 실패 1")
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

	// device 도메인. 인증이 성공했을때
	if userAuthResult {

		deviceRegistStateInput := adapter.MakeDeviceRegistStateInput(req.Id, req.Uuid)
		challenge, err := h.svc.Device.GetDeviceRegistState(ctx, deviceRegistStateInput)

		if err != nil {
			if err == consts.ErrDeviceNotRegist {
				log.Println("[UserAuthAndDeviceCheck] device 등록을 위한 challenge 발급")
				// 등록되지 않은 device이므로 challenge 발급함.
				res := userAuthService.UserAuthServiceResponse{
					DeviceChallenge: challenge,
				}
				response.SendSuccess(c, res)
			} else {
				log.Println("[UserAuthAndDeviceCheck] 인증 실패 2")
				response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
			}
			return
		}

		// device는 등록 완료.
		// token 도메인. at, rt 발급 체크 및 response
		userAuthTokenInput := input.MakeGenerateAuthTokenInput(req.Id, req.Uuid)
		output, err := h.svc.Token.GenerateAuthToken(ctx, userAuthTokenInput)
		if err != nil {
			log.Println("[UserAuthAndDeviceCheck] 인증 실패 3")
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
			return
		}

		res := userAuthService.UserAuthServiceResponse{
			RefreshToken:    output.RefreshToken,
			RefreshTokenExp: output.RefreshTokenExp,
			AccessToken:     output.AccessToken,
		}
		response.SendSuccess(c, res)
		return

	} else {
		// 인증 실패.
		log.Println("[UserAuthAndDeviceCheck] 인증 실패 4")
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}
}
