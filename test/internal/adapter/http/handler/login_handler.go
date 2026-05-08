package handler

import (
	"context"
	"encoding/json"
	"test/internal/adapter/http/dto/login"
	"test/internal/adapter/http/mapper"
	"test/internal/application/usecase"
	commonConsts "test/internal/consts"
	"test/pkg/response"
	"time"

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

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()하게되면 go routine 내부로직이 수행되기도 전에 끝나버림. 그래서 login_usecase.go의 login 함수가 <-ctx.Done()에서 걸릴 틈도 없이 cancel이 수행되버리는 구조.

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
