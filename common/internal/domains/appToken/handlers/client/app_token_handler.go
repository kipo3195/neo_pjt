package client

import (
	"common/internal/consts"
	requestDTO "common/internal/domains/appToken/dto/client/requestDTO"
	usecases "common/internal/domains/appToken/usecases/client"
	"encoding/json"
	"errors"

	commonConsts "common/pkg/consts"
	"common/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AppTokenHandler struct {
	usecase usecases.AppTokenUsecase
}

func NewAppTokenHandler(usecase usecases.AppTokenUsecase) *AppTokenHandler {
	return &AppTokenHandler{
		usecase: usecase,
	}
}

func (h *AppTokenHandler) AppTokenRefresh(c *gin.Context) {

	ctx := c.Request.Context()

	// request body 데이터 -> dto로 변경
	var body requestDTO.AppTokenRefreshRequestBody

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

	requestDTO := requestDTO.AppTokenRefreshRequestDTO{
		Body: body,
	}

	data, err := h.usecase.AppTokenReIssueInAuth(ctx, requestDTO)

	if err != nil {
		switch {
		case errors.Is(err, consts.ErrRefreshTokenAuthInvalid):
			// 토큰 검증 실패
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F003, consts.COMMON_F003_MSG)
		case errors.Is(err, consts.ErrRefreshTokenAuthExpired):
			// 토큰 만료
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_F004, consts.COMMON_F004_MSG)
		default:
			// 서버 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	} else {
		response.SendSuccess(c, data)
	}
}
