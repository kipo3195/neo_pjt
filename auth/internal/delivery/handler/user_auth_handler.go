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

	ctx := c.Request.Context()

	var req userAuth.UserAuthChallengeRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	userAuthChallengeInput := input.MakeUserAuthChallengeInput(req.Id, req.Device)
	userAuthChallengeOutput, err := h.usecase.GenerateUserAuthChallenge(ctx, userAuthChallengeInput)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, userAuthChallengeOutput)
	}

}

func (h UserAuthHandler) GetUserAuth(c *gin.Context) {

	ctx := c.Request.Context()
	var req userAuth.UserAuthRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	userAuthInput := input.MakeUserAuthInput(req.Id, req.Fv, req.Device)
	userAuthOutput := h.usecase.GetUserAuth(ctx, userAuthInput)

	res := userAuth.UserAuthResponse{
		AccessToken:     userAuthOutput.AccessToken,
		RefreshToken:    userAuthOutput.RefreshToken,
		DeviceChallenge: userAuthOutput.DeviceChallenge,
	}

	response.SendSuccess(c, res)

}

func (h UserAuthHandler) UserAuthInfoRegister(c *gin.Context) {

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
