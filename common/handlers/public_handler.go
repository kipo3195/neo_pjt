package handlers

import (
	clCommonReqDto "common/dto/client/request"
	"common/usecases"
	"common/utils"
	"encoding/json"
	"errors"
	"net/http"

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

func (h *PublicHandler) AppValidation(c *gin.Context) {

	ctx := c.Request.Context()

	// request body 데이터 -> dto로 변경
	body := clCommonReqDto.AppValidationRequestBody{
		Uuid:       c.Query("uuid"),
		AppToken:   c.Query("appToken"),
		Device:     c.Query("device"),
		SkinHash:   c.Query("skinHash"),
		ConfigHash: c.Query("configHash"),
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clCommonReqDto.AppValidationRequestDTO{
		Body: body,
	}

	// 검증
	data, err := h.usecase.AppValidation(ctx, requestDTO)

	if err != nil || !data {
		switch {
		case errors.Is(err, consts.ErrInvalidClaims):
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_106, consts.E_106_MSG)
		case errors.Is(err, consts.ErrSkinHashInvalid):
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.COMMON_F001, consts.COMMON_F001_MSG)
		case errors.Is(err, consts.ErrConfigHashInvalid):
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.COMMON_F002, consts.COMMON_F002_MSG)
		default:
			utils.SendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}
		return
	} else if data {
		utils.SendSuccessResponse(c, "")
	}
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
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.COMMON_F003, consts.COMMON_F003_MSG)
			writeJSON(w, http.StatusBadRequest, res)
		case errors.Is(err, consts.ErrRefreshTokenAuthExpired):
			// 토큰 만료
			res.Result = consts.FAIL
			res.Data = newErrorResp(consts.COMMON_F004, consts.COMMON_F004_MSG)
			writeJSON(w, http.StatusBadRequest, res)
		default:
			// 서버 에러
			res.Result = consts.ERROR
			res.Data = newErrorResp(consts.E_500, consts.E_500_MSG)
			writeJSON(w, http.StatusInternalServerError, res)
		}
		return
	}

	res.Result = consts.SUCCESS
	res.Data = data
	writeJSON(w, http.StatusOK, res)
}
