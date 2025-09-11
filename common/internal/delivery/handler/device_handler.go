package handler

import (
	"common/internal/application/orchestrator"
	"common/internal/delivery/adapter"
	"common/internal/delivery/dto/device"
	"common/pkg/response"
	"encoding/json"

	commonConsts "common/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type DeviceHandler struct {
	svc *orchestrator.DeviceInitService
}

func NewDeviceHandler(svc *orchestrator.DeviceInitService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) DeviceInit(c *gin.Context) {

	// request body 데이터 -> dto로 변경
	var body device.DeviceRequestBody
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

	requestDTO := device.DeviceDTO{
		Body: body,
	}

	deviceWrapper := adapter.DeviceWrapper{DeviceRequestBody: requestDTO.Body}
	connectInfoInput := adapter.MakeConnectInfoInput(deviceWrapper)

	worksInfo, err := h.svc.Device.GetConnectInfo(connectInfoInput)
	// 호출 도메인(DNS)만 뽑아내는데 도메인 명칭이 Device?  차라리 configuration에서 통합 관리 또는 worksInfo 도메인을 별도로 생성한다면?
	// 서버 정보 (도메인), worksCode, worksName, useYn, regDate
	// 그리고 값을 매번 DB에서 조회할 필요가 있나?
	if err != nil {

	}

	var serverUrl string // 호출 도메인을 주입받아서 사용할 수 있도록 하기
	// 수정완료

	issuedAppToken, err := h.svc.AppToken.GenerateAppTokenInAuth(body.Uuid, serverUrl)
	if err != nil {

	}

	// 수정완료
	worksConfig := h.svc.Configuration.GetWorksConfig()
	if err != nil {

	}

	// 수정완료
	skinInfo, err := h.svc.Skin.GetSkinInfo()
	if err != nil {

	}

	// 결국 수정되어야할 api의 방향
	result := device.DeviceResultResponse{
		WorksInfo:      worksInfo, // works의 정보
		IssuedAppToken: issuedAppToken,
		WorksConfig:    worksConfig, // works의 설정정보
		SkinInfo:       skinInfo,
	}

	response.SendSuccess(c, result)
}
