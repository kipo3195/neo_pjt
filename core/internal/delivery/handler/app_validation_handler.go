package handler

import (
	"core/internal/application/usecase"
	"core/internal/consts"
	"core/internal/delivery/adapter"
	"core/internal/delivery/dto/appValidation"
	"core/internal/infrastructure/config"
	commonConsts "core/pkg/consts"
	"encoding/json"
	"errors"

	response "core/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppValidationHandler struct {
	usecase usecase.AppValidationUsecase
	sfg     *config.ServerConfig
}

func NewAppValidationHandler(sfg *config.ServerConfig, uc usecase.AppValidationUsecase) *AppValidationHandler {
	return &AppValidationHandler{usecase: uc, sfg: sfg}
}

func (h *AppValidationHandler) ValidateApp(c *gin.Context) {

	// request의 header 데이터 -> dto로 변경
	var headerPrefix = consts.HEADER_PRIFIX
	header := appValidation.AppValidationRequestHeader{
		Hash:   c.GetHeader(headerPrefix + "Hash"),
		Device: c.GetHeader(headerPrefix + "Device"),
		Uuid:   c.GetHeader(headerPrefix + "Uuid"),
	}

	// header 검증
	if header.Hash == "" || header.Device == "" || header.Uuid == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body appValidation.AppValidationRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	requestDTO := appValidation.AppValidationRequestDTO{
		Header: header,
		Body:   body,
	}

	// adapter : dto <-> input
	validationInput := adapter.MakeValidateAppInput(requestDTO)

	// 배포 앱 hash 검증
	result, err := h.usecase.CheckValidation(c, validationInput)

	validationOutput := adapter.MakeValidateAppOutput(result)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 조회된 결과 없음 - 앱 해시 검증 실패
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.CORE_F101, consts.CORE_F101_MSG)

		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}
	response.SendSuccess(c, validationOutput)

	// 클라이언트가 넘겨준 Domain : 테넌트 정보로 검증

	// resDto, err := h.usecase.GetWorksInfos(validationInput)

	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		// 조회된 결과 없음 - 매핑된 서버 정보가 없을때 에러
	// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.CORE_F102, consts.CORE_F102_MSG)
	// 	} else {
	// 		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	// 	}
	// 	return
	// } else {
	// 	response.SendSuccess(c, resDto.Body)
	// }
}
