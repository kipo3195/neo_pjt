package server

import (
	"common/internal/consts"
	"common/internal/domains/appValidation/dto/server/requestDTO"
	usecases "common/internal/domains/appValidation/usecases/server"
	"errors"

	commonConsts "common/pkg/consts"
	response "common/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AppValidationHandler struct {
	usecase usecases.AppValidationUsecase
}

func NewAppValidationHandler(usecase usecases.AppValidationUsecase) *AppValidationHandler {

	return &AppValidationHandler{
		usecase: usecase,
	}
}

func (h *AppValidationHandler) AppValidation(c *gin.Context) {

	ctx := c.Request.Context()

	// request body 데이터 -> dto로 변경
	body := requestDTO.AppValidationRequestBody{
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

	requestDTO := requestDTO.AppValidationRequestDTO{
		Body: body,
	}

	// 검증
	data, err := h.usecase.AppValidation(ctx, requestDTO)

	if err != nil || !data {
		switch {
		case errors.Is(err, consts.ErrInvalidClaims):
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_106, commonConsts.E_106_MSG)
		case errors.Is(err, consts.ErrSkinHashInvalid):
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F001, consts.COMMON_F001_MSG)
		case errors.Is(err, consts.ErrConfigHashInvalid):
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F002, consts.COMMON_F002_MSG)
		default:
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	} else if data {
		response.SendSuccess(c, "")
	}
}
