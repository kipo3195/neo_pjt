package handler

import (
	"admin/internal/application/orchestrator"
	"admin/internal/delivery/adapter"
	"admin/internal/delivery/dto/serviceUser"
	"admin/pkg/consts"
	response "admin/pkg/response"
	"encoding/json"

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

func (r *ServiceUserAuthRegisterHandler) RegistServiceUser(c *gin.context) {

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

	input := adapter.MakeRegistServiceUserInput(req.Org, req.UserId, req.UserAuth)
	output, err := r.svc.ServiceUser.RegistServiecUser(ctx, input)

	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		return
	}

	user := make([]serviceUser.ServiceUser, 0)

	for _, s := range output.ServiceUser {

		value := serviceUser.ServiceUser{
			UserId: s.UserId,
			Hash:   s.UserHash,
		}

		user = append(user, value)
	}

	res := serviceUser.RegistServiceUserResponse{
		ServiceUser: user,
	}

	response.SendSuccess(c, res)

}
