package handler

import (
	"auth/internal/application/usecase"
	"auth/internal/application/usecase/input"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"

	"auth/internal/delivery/dto/userAuth"

	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	usecase usecase.UserAuthUsecase
}

func NewUserAuthHandler(uc usecase.UserAuthUsecase) *UserAuthHandler {
	return &UserAuthHandler{
		usecase: uc,
	}
}

func (h UserAuthHandler) GenerateAuthChallenge(c *gin.Context) {

}

func (h UserAuthHandler) GetAuthStatus(c *gin.Context) {

}

func (h UserAuthHandler) UserAuthRegister(c *gin.Context) {

	ctx := c.Request.Context()

	var req userAuth.UserAuthRegisterRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	userAuthRegisterInput := input.MakeUserAuthRegisterInput(req.Id, req.Salt, req.AuthHash, req.UserHash)
	userAuthRegisterOutput := h.usecase.PutUserAuth(ctx, userAuthRegisterInput)

	response.SendSuccess(c, userAuthRegisterOutput)

}
