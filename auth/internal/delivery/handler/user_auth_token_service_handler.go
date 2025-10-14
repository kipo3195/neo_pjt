package handler

import (
	"auth/internal/application/orchestrator"
	"auth/internal/application/usecase/input"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"

	"auth/internal/delivery/adapter"
	"auth/internal/delivery/dto/userAuthTokenService"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserAuthTokenServiceHandler struct {
	svc *orchestrator.UserAuthTokenService
}

func NewUserAuthTokenServiceHandler(svc *orchestrator.UserAuthTokenService) *UserAuthTokenServiceHandler {
	return &UserAuthTokenServiceHandler{
		svc: svc,
	}
}

func (h *UserAuthTokenServiceHandler) AccessTokenRefresh(c *gin.Context) {

	ctx := c.Request.Context()
	var req userAuthTokenService.AccessTokenReIssueRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	// appToken 검증
	appTokenValidationInput := adapter.MakeAppTokenValidationInput(req.AppToken, req.Token, req.TokenType, req.Uuid)
	result, err := h.svc.Token.AppTokenValidation(appTokenValidationInput, ctx)

	// accessToken 재발급
	if result {
		accessTokenReIssueInput := input.MakeAccessTokenReIssueInput(req.AppToken, req.Uuid, req.RefreshToken)
		output, err := h.svc.UserAuth.AccessTokenReIssue(ctx, accessTokenReIssueInput)

		if
	} else if err != nil {
		// error type custom todo
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		// error type custom todo
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}
