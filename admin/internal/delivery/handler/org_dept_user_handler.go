package handler

import (
	"admin/internal/application/usecase"
	"admin/internal/application/usecase/input"
	"admin/internal/delivery/dto/orgDeptUser"
	"admin/pkg/consts"
	"encoding/json"

	response "admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgDeptUserHandler struct {
	usecase usecase.OrgDeptUsersUsecase
}

func NewOrgDeptUsersHandler(usecase usecase.OrgDeptUsersUsecase) *OrgDeptUserHandler {
	return &OrgDeptUserHandler{
		usecase: usecase,
	}

}

func (h *OrgDeptUserHandler) RegistDeptUser(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req orgDeptUser.RegistDeptUserRequest

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

	input := input.MakeRegistDeptUserInput(req.UserHash, req.DeptCode, req.DeptOrg, req.PositionCode, req.RoleCode, req.IsConcurrentPosition)

	// usecase 호출
	_, err := h.usecase.RegistDeptUser(ctx, input)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *OrgDeptUserHandler) DeleteDeptUser(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req orgDeptUser.DeleteDeptUserRequest

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

	input := input.MakeDeleteDeptUserInput(req.UserHash, req.DeptCode, req.DeptOrg)

	// usecase 호출
	_, err := h.usecase.DeleteDeptUser(ctx, input)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *OrgDeptUserHandler) GetDeptUser(c *gin.Context) {

}

func (h *OrgDeptUserHandler) UpdateDeptUser(c *gin.Context) {

}
