package handlers

import (
	"core/config"
	consts "core/consts"
	clReqDto "core/dto/client/request"
	"core/usecases"
	"core/utils"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CoreHandler struct {
	usecase usecases.CoreUsecase
	sfg     *config.ServerConfig
}

func NewCoreHandler(sfg *config.ServerConfig, uc usecases.CoreUsecase) *CoreHandler {
	return &CoreHandler{usecase: uc, sfg: sfg}
}

// /app-validation
func (h *CoreHandler) GetAppValidation(c *gin.Context) {

	// request의 header 데이터 -> dto로 변경
	var headerPrefix = h.sfg.ApiConfig.NeoHeaderPrefix
	header := clReqDto.AppValidationRequestHeader{
		Hash:   c.GetHeader(headerPrefix + "Hash"),
		Device: c.GetHeader(headerPrefix + "Device"),
		Uuid:   c.GetHeader(headerPrefix + "Uuid"),
	}

	// header 검증
	if header.Hash == "" || header.Device == "" || header.Uuid == "" {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body clReqDto.AppValidationRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	requestDTO := clReqDto.AppValidationRequestDTO{
		Header: header,
		Body:   body,
	}

	// 배포 앱 hash 검증
	_, err := h.usecase.CheckValidation(requestDTO.Header)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 조회된 결과 없음 - 앱 해시 검증 실패
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.CORE_F101, consts.CORE_F101_MSG)
		} else {
			utils.SendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}
		return
	}

	// 클라이언트가 넘겨준 Domain : 테넌트 정보로 검증
	resDto, err := h.usecase.GetWorksInfos(requestDTO)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 조회된 결과 없음 - 매핑된 서버 정보가 없을때 에러
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.CORE_F102, consts.CORE_F102_MSG)
		} else {
			utils.SendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}
		return
	} else {
		utils.SendSuccessResponse(c, resDto.Body)
	}
}
