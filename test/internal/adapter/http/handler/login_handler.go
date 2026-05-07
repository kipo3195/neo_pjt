package handler

import (
	"encoding/json"
	"test/internal/adapter/http/dto/login"
	"test/internal/adapter/http/mapper"
	"test/internal/application/usecase"
	commonConsts "test/internal/consts"
	"test/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	usecase usecase.LoginUsecase
}

func NewLoginHandler(usecase usecase.LoginUsecase) *LoginHandler {
	return &LoginHandler{
		usecase: usecase,
	}
}

func (r *LoginHandler) PutLogin(c *gin.Context) {

	ctx := c.Request.Context()

	var req login.LoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	input := mapper.MakePutLoginInput(req.UserIdPrefix, req.ConnectionCount)
	// 별도의 고루틴으로 처리
	go r.usecase.PutLogin(ctx, input)

	response.SendSuccess(c, "success")
}
