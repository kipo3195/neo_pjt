package handler

import (
	"admin/internal/application/usecase"
	"admin/internal/application/usecase/input"
	"admin/internal/delivery/dto/orgDept"
	"admin/pkg/consts"
	"encoding/json"

	response "admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgDeptHandler struct {
	usecase usecase.OrgDeptsUsecase
}

func NewOrgDeptsHandler(usecase usecase.OrgDeptsUsecase) *OrgDeptHandler {
	return &OrgDeptHandler{
		usecase: usecase,
	}

}

func (h *OrgDeptHandler) RegisterDept(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	// 이미 timeout이 걸려있는 ctx이므로 그대로 사용만 하면됨.
	ctx := c.Request.Context()

	var req orgDept.RegisterDeptRequest
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

	// 필요시 result 값 response status로 분기 처리
	input := input.MakeRegisterDeptInput(req)
	_, err := h.usecase.RegisterDept(ctx, input)

	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}
}

func (h *OrgDeptHandler) GetDept(c *gin.Context) {

}

func (h *OrgDeptHandler) UpdateDept(c *gin.Context) {

}

func (h *OrgDeptHandler) DeleteDept(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req orgDept.DeleteDeptRequest

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

	input := input.MakeDeleteDeptInput(req.DeptOrg, req.DeptCode)

	// usecase 호출
	_, err := h.usecase.DeleteDept(ctx, input)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}
