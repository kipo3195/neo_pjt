package handler

import (
	"auth/internal/application/usecase"
	"auth/internal/consts"
	"auth/internal/delivery/adapter"
	"auth/internal/delivery/dto/device"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(uc usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: uc}
}

func (h *DeviceHandler) GetMyDeviceInfo(c *gin.Context) {

	ctx := c.Request.Context()

	hash := c.Value(consts.USER_HASH)
	myHash, ok := hash.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	myDeviceInput := adapter.MakeMyDeviceInfoInput(myHash)
	output, err := h.usecase.GetMyDeviceInfo(ctx, myDeviceInput)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	deviceInfoArr := make([]device.DeviceInfo, 0)

	for _, temp := range output {

		deviceInfo := device.DeviceInfo{
			Uuid:      temp.Uuid,
			Version:   temp.Version,
			ModelName: temp.ModelName,
			CreateAt:  temp.CreateAt,
		}
		deviceInfoArr = append(deviceInfoArr, deviceInfo)
	}

	res := device.GetMyDeviceInfoResponse{
		MyDeviceInfo: deviceInfoArr,
	}

	response.SendSuccess(c, res)
}
