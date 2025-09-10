package handler

import (
	"common/internal/application/orchestrator"
	"common/internal/consts"
	"common/pkg/response"

	commonConsts "common/pkg/consts"

	"common/internal/delivery/dto/appValidation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AppValidationHandler struct {
	svc *orchestrator.AppValidationService
}

func NewAppValidationHander(svc *orchestrator.AppValidationService) *AppValidationHandler {
	return &AppValidationHandler{svc: svc}
}

// POST /server/v1/app-validation
func (h *AppValidationHandler) GetAppValidation(c *gin.Context) {

	// 실제 비즈니스 로직 처리? svc를 호출 기존 handler와 동일하게 처리하도록 수정 필요.
	body := appValidation.AppValidationRequestBody{
		Uuid:        c.Query("uuid"),
		AppToken:    c.Query("appToken"),
		AccessToken: c.Query("accesssToken"),
		Device:      c.Query("device"),
		SkinHash:    c.Query("skinHash"),
		ConfigHash:  c.Query("configHash"),
	}
	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	requestDTO := appValidation.AppValidationRequestDTO{
		Body: body,
	}

	domainReq := toAppValidationDomainDTO(requestDTO.Body)

	// 기존 appValidtion 도메인에서 하던 처리
	data, err := h.svc.Validator.AppValidation(c, domainReq)

	if err != nil || !data {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 기존 skin 에서 하던 처리
	_, err = h.svc.Skin.CheckSkin(domainReq.Body.SkinHash)
	if err != nil {
		// 스킨 정보 에러로 변경
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F001, consts.COMMON_F001_MSG)
		return
	}

	// 기존 config에서 하던 처리
	_, err = h.svc.Configuration.CheckConfiguration(domainReq.Body.ConfigHash)
	if err != nil {
		// 스킨 정보 에러로 변경
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F002, consts.COMMON_F002_MSG)
		return
	}

	response.SendSuccess(c, "")
}

// 변환
func toAppValidationDomainDTO(serviceBody appValidation.AppValidationRequestBody) appValidation.AppValidationRequestDTO {
	return appValidation.AppValidationRequestDTO{
		Body: appValidation.AppValidationRequestBody{
			Uuid:        serviceBody.Uuid,
			AppToken:    serviceBody.AppToken,
			AccessToken: serviceBody.AccessToken,
			Device:      serviceBody.Device,
			SkinHash:    serviceBody.SkinHash,
			ConfigHash:  serviceBody.ConfigHash,
		},
	}
}
