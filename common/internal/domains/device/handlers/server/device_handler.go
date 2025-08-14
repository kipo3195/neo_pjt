package server

import (
	"common/internal/consts"
	"common/internal/domains/device/dto/server/requestDTO"
	deviceUsecase "common/internal/domains/device/usecases/server"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	usecase deviceUsecase.DeviceUsecase
}

func NewDeviceHandler(usecase deviceUsecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{
		usecase: usecase,
	}
}

func (h *DeviceHandler) DeviceInit(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request body 데이터 -> dto로 변경
	var body *requestDTO.DeviceInitRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	fmt.Printf("CORE 서버에서 호출, uuid : %s, worksCode : %s \n", body.Uuid, body.WorksCode)

	// DB에서 해당 works의 정보조회 + AUTH에서 토큰 발급 요청
	data, err := h.usecase.DeviceInit(ctx, body)

	if err != nil {
		if err == consts.ErrDB {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_102, commonConsts.E_102_MSG)
		} else {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
	} else {
		response.SendSuccess(c, data)
	}
}
