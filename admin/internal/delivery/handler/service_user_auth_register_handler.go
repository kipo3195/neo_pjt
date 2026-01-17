package handler

import (
	"admin/internal/application/orchestrator"
	"admin/internal/delivery/adapter"
	"admin/internal/delivery/dto/serviceUser"
	"admin/pkg/consts"
	response "admin/pkg/response"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ServiceUserAuthRegisterHandler struct {
	svc *orchestrator.ServiceUserAuthRegisterService
}

func NewServiceUserAuthRegisterHandler(svc *orchestrator.ServiceUserAuthRegisterService) *ServiceUserAuthRegisterHandler {
	return &ServiceUserAuthRegisterHandler{
		svc: svc,
	}
}

func (r *ServiceUserAuthRegisterHandler) RegistServiceUser(c *gin.Context) {

	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req serviceUser.RegistServiceUserRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	serviceUserInput := adapter.MakeRegistServiceUserInput(req.Org, req.UserId, req.UserAuth)
	serviceUserOutput, err := r.svc.ServiceUser.RegistServiecUser(ctx, serviceUserInput)

	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		return
	}

	userAuthInput := adapter.MakeUserAuthInput(serviceUserOutput)
	err = r.svc.UserAuthRegister.UserAuthRegisterInAuth(ctx, userAuthInput)

	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		return
	}

	publishServiceUserInput := adapter.MakePublishServiceUserInput(req.Org, serviceUserOutput)
	err = r.svc.ServiceUser.PublishServiceUser(ctx, publishServiceUserInput)

	user := make([]serviceUser.ServiceUser, 0)

	for _, s := range serviceUserOutput.ServiceUser {

		value := serviceUser.ServiceUser{
			UserId: s.UserId,
			Hash:   s.UserHash,
			// Salt, UserAuth는 정의하지 않았음.
		}

		user = append(user, value)
	}

	res := serviceUser.RegistServiceUserResponse{
		ServiceUser: user,
	}

	response.SendSuccess(c, res)

}
