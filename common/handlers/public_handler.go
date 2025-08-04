package handlers

import (
	clCommonReqDto "common/dto/client/request"
	"common/usecases"
	"common/utils"
	"encoding/json"
	"errors"

	consts "common/consts"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type PublicHandler struct {
	usecase usecases.PublicUsecase
}

func NewPublicHandler(uc usecases.PublicUsecase) *PublicHandler {
	return &PublicHandler{usecase: uc}
}

func (h *PublicHandler) AppTokenRefresh(c *gin.Context) {

	ctx := c.Request.Context()

	// request body 데이터 -> dto로 변경
	var body clCommonReqDto.AppTokenRefreshRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clCommonReqDto.AppTokenRefreshRequestDTO{
		Body: body,
	}

	data, err := h.usecase.AppTokenReIssue(ctx, requestDTO)

	if err != nil {
		switch {
		case errors.Is(err, consts.ErrRefreshTokenAuthInvalid):
			// 토큰 검증 실패
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.COMMON_F003, consts.COMMON_F003_MSG)
		case errors.Is(err, consts.ErrRefreshTokenAuthExpired):
			// 토큰 만료
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.COMMON_F004, consts.COMMON_F004_MSG)
		default:
			// 서버 에러
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}
		return
	} else {
		utils.SendSuccessResponse(c, data)
	}
}
