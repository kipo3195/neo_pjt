package handlers

import (
	domain "common/internal/domains/device/dto/server/requestDTO"
	serviceDto "common/internal/serviceDto"
	"common/internal/services"
	"common/pkg/response"
	"encoding/json"

	commonConsts "common/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type DeviceInitHandler struct {
	svc *services.DeviceInitService
}

func NewDeviceInitHandler(svc *services.DeviceInitService) *DeviceInitHandler {
	return &DeviceInitHandler{svc: svc}
}

func (h *DeviceInitHandler) DeviceInit(c *gin.Context) {

	// request body 데이터 -> dto로 변경
	var body serviceDto.DeviceInitRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	connectInfo, err := h.svc.Device.GetConnectInfo(toDeviceDomainDTO(body))
	if err != nil {

	}

	issuedAppToken, err := h.svc.Device.GetIssuedAppToken(toDeviceDomainDTO(body))

	if err != nil {

	}

	worksConfig, err := h.svc.Device.GetWorksConfig(toDeviceDomainDTO(body))
	if err != nil {

	}

	result := serviceDto.DeviceInitResultResponse{
		ConnectInfo:    connectInfo,
		IssuedAppToken: issuedAppToken,
		WorksConfig:    worksConfig,
	}

	response.SendSuccess(c, result)
}

func toDeviceDomainDTO(body serviceDto.DeviceInitRequestBody) *domain.DeviceInitRequest {
	return &domain.DeviceInitRequest{
		WorksCode: body.WorksCode,
		Uuid:      body.Uuid,
		Device:    body.Device,
	}
}
