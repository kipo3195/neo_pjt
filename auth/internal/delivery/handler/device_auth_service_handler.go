package handler

import (
	"auth/internal/application/orchestrator"
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

	// 20251125 채팅, 쪽지 내용 암호화 공개키 암호화 처리 후 저장 (chKey, noKey)
	// message service API 호출
	otpKeyRegistInput := adapter.MakeOtpKeyRegistInput(req.Id, req.Uuid, req.ChKey, req.NoKey)
	otpKeyRegistResult, err := h.svc.Otp.OtpKeyRegist(ctx, otpKeyRegistInput)

	if err != nil {

		return
	}

	// 디바이스 정보 등록 체크
	deviceRegistInput := adapter.MakeDeviceRegistCheckInput(req.Id, req.Uuid, req.Challenge)
	deviceRegResult, err := h.svc.Device.DeviceRegistCheck(ctx, deviceRegistInput)

	if err != nil {
		// 등록에 따라 다르게 처리 필요 TODO challege 만료
		if err == consts.ErrDeviceChallengeExpired {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F007, consts.AUTH_F007_MSG)

			// challenge 불일치
		} else if err == consts.ErrDeviceChallengeMismatch {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F012, consts.AUTH_F012_MSG)

			// 이미 등록된 device
		} else if err == consts.ErrDeviceAlreadyRegist {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F013, consts.AUTH_F013_MSG)

		} else {
			// 저장실패 DB error
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}

		return
	}

	if deviceRegResult {
		// 등록 성공
		generateAuthTokenInput := adapter.MakeGenerateAuthTokenInput(req.Id, req.Uuid, true)
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
			ChKeyRegDate:    otpKeyRegistResult.ChKeyRegDate,
			NoKeyRegDate:    otpKeyRegistResult.NoKeyRegDate,
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

	// 가장 최신의 rt로 요청했는지 점검
	refreshTokenCheckInput := adapter.MakeRefreshTokenCheckInput(req.Uuid, req.RefreshToken)
	id, err := h.svc.Token.CheckRefreshToken(refreshTokenCheckInput, ctx)

	if err != nil {
		if err == consts.ErrUserIdDoesNotExist {
			// uuid : refresh token에 매핑된 사용자 id가 없을때
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, commonConsts.E_109, commonConsts.E_109_MSG)
		} else if err == consts.ErrRefreshTokenAuthError {
			// 가장 최신의 RT로 요청하지 않음.
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F011, consts.AUTH_F011_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	// device 등록 점검 (재발급 받을 수 있는지)
	deviceRegistInput := adapter.MakeDeviceRegistStateInput(id, req.Uuid)
	challenge, err := h.svc.Device.GetDeviceRegistState(ctx, deviceRegistInput)

	log.Println("[DeviceRefresh] id : ", id, ", new challenge : ", challenge)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// device 정보 업데이트
	updateDeviceInfoInput := adapter.MakeUpdateDeviceInfoInput(id, req.Uuid, req.ModelName, req.Version)
	err = h.svc.Device.UpdateDeviceInfo(ctx, updateDeviceInfoInput)

	if err != nil {
		if err == consts.ErrDeviceNotRegist {
			// 디바이스 등록되있지도 않은데 재발급 요청함
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F010, consts.AUTH_F010_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	// at, rt 신규 생성 - 무조건 신규 발급. 만료되었을 가능성이 높으니..
	generateAuthTokenInput := adapter.MakeGenerateAuthTokenInput(id, req.Uuid, true)
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
	removeDeviceChallengeInput := adapter.MakeRemoveDeviceChallengeInput(id, req.Uuid)
	h.svc.Device.RemoveDeviceChallenge(ctx, removeDeviceChallengeInput)

	response.SendSuccess(c, res)

}
