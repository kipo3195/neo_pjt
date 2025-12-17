package handler

import (
	"auth/internal/application/usecase"
	"auth/internal/consts"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
	"encoding/json"
	"errors"
	"fmt"

	"auth/internal/delivery/adapter"
	"auth/internal/delivery/dto/userAuth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	userAuthChallengeInput := adapter.MakeUserAuthChallengeInput(req.Id, req.Uuid)
	temp, err := h.usecase.GenerateUserAuthChallenge(ctx, userAuthChallengeInput)

	if err != nil {
		if errors.Is(err, consts.ErrSaltNotRegist) {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, consts.AUTH_F006, consts.AUTH_F006_MSG)
		} else {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	userAuthChallengeOutput := adapter.MakeUserAuthChallengeOutput(temp.Challenge, temp.Salt)

	res := userAuth.UserAuthChallengeResponse{
		Challenge: userAuthChallengeOutput.Challenge,
		Salt:      userAuthChallengeOutput.Salt,
	}

	response.SendSuccess(c, res)

}

func (h UserAuthHandler) GetUserAuth(c *gin.Context) {

	ctx := c.Request.Context()
	var req userAuth.UserAuthRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		fmt.Println(err)
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	userAuthInput := adapter.MakeUserAuthInput(req.Id, req.Fv, req.Uuid)
	userAuthOutput, err := h.usecase.GetUserAuth(ctx, userAuthInput)

	if err != nil {
		if errors.Is(err, consts.ErrUserAuthFvMismatch) {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
		} else if errors.Is(err, consts.ErrUserAuthChallengeExpired) {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F007, consts.AUTH_F007_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	if userAuthOutput {
		res := userAuth.UserAuthResponse{
			AccessToken:     "",
			RefreshToken:    "",
			DeviceChallenge: "",
		}
		response.SendSuccess(c, res)
	}

}

func (h UserAuthHandler) UserAuthInfoRegister(c *gin.Context) {

	ctx := c.Request.Context()

	var req userAuth.UserAuthRegisterRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	userAuthRegisterInput := adapter.MakeUserAuthRegisterInput(req.UserAuth)
	userAuthRegisterOutput := h.usecase.PutUserAuth(ctx, userAuthRegisterInput)

	response.SendSuccess(c, userAuthRegisterOutput)

}
