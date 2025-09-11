package handler

import (
	"common/internal/application/usecase"
	"common/internal/application/usecase/input"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"encoding/json"

	"common/internal/delivery/dto/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) UserRegister(c *gin.Context) {

	ctx := c.Request.Context()

	var req user.UserRegisterRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	userRegisterInput := input.MakeUserRegisterInput(req.Id, req.Salt, req.Fv)
	userRegisterOutput := h.usecase.UserRegister(ctx, userRegisterInput)

	if userRegisterOutput == "code1" {
		response.SendSuccess(c, "success")
	} else if userRegisterOutput == "code2" {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.FAIL, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *UserHandler) GetUserRegisterChallenge(c *gin.Context) {

	ctx := c.Request.Context()
	var req = user.UserRegisterChallengeRequest{
		Id:   c.Query("id"),
		Salt: c.Query("salt"),
	}

	// 유효성 검증 로직
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	userRegisterChallengeInput := input.MakeUserRegisterChallengeInput(req.Id, req.Salt)

	userRegisterChallengeOutput := h.usecase.GenerateUserChallenge(ctx, userRegisterChallengeInput)

	// output -> dto ?
	res := user.UserRegisterChallengeResponse{
		Challenge: userRegisterChallengeOutput.Challenge,
	}

	response.SendSuccess(c, res)

}
