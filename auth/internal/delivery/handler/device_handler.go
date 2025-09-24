package handler

import (
	"auth/internal/application/usecase"
	"auth/internal/application/usecase/input"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"
	"fmt"

	"auth/internal/delivery/dto/device"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(uc usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: uc}
}

func (h DeviceHandler) DeviceRegist(c *gin.Context) {

	ctx := c.Request.Context()
	var req device.DeviceRegistRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		fmt.Println(err)
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	deviceRegistInput := input.MakeDeviceRegistInput(req.Id, req.Uuid, req.ModelName, req.Version, req.Challenge)
	deviceRegistOutput, err := h.usecase.DeviceRegist(ctx, deviceRegistInput)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}
	res := device.DeviceRegistResponse{
		AccessToken:  deviceRegistOutput.AccessToken,
		RefreshToken: deviceRegistOutput.RefreshToken,
	}

	response.SendSuccess(c, res)
}
