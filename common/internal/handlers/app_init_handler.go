package handlers

import (
	"common/internal/consts"
	"common/internal/services"
	"common/pkg/response"

	commonConsts "common/pkg/consts"

	domain "common/internal/domains/appValidation/dto/server/requestDTO"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AppInitHandler struct {
	svc *services.AppInitService
}

func NewAppInitHander(svc *services.AppInitService) *AppInitHandler {
	return &AppInitHandler{svc: svc}
}

type AppValidationRequestDTO struct {
	Body AppValidationRequestBody
}

type AppValidationRequestBody struct {
	AppToken   string `json:"appToken" validate:"required"`
	Uuid       string `json:"uuid" validate:"required"`
	Device     string `json:"device" validate:"required"`
	SkinHash   string `json:"skinHash" validate:"required"`
	ConfigHash string `json:"configHash" validate:"required"`
}

// POST /server/v1/app-validation
func (h *AppInitHandler) GetAppValidation(c *gin.Context) {

	// 실제 비즈니스 로직 처리? svc를 호출 기존 handler와 동일하게 처리하도록 수정 필요.
	body := AppValidationRequestBody{
		Uuid:       c.Query("uuid"),
		AppToken:   c.Query("appToken"),
		Device:     c.Query("device"),
		SkinHash:   c.Query("skinHash"),
		ConfigHash: c.Query("configHash"),
	}
	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	requestDTO := AppValidationRequestDTO{
		Body: body,
	}

	domainReq := toDomainDTO(requestDTO.Body)

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
func toDomainDTO(serviceBody AppValidationRequestBody) domain.AppValidationRequestDTO {
	return domain.AppValidationRequestDTO{
		Body: domain.AppValidationRequestBody{
			Uuid:       serviceBody.Uuid,
			AppToken:   serviceBody.AppToken,
			Device:     serviceBody.Device,
			SkinHash:   serviceBody.SkinHash,
			ConfigHash: serviceBody.ConfigHash,
		},
	}
}
