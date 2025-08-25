package server

import (
	commonConsts "org/pkg/consts"

	"encoding/json"
	"org/internal/domains/department/dto/server/requestDTO"
	usecases "org/internal/domains/department/usecases/server"
	"org/pkg/response"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	usecase usecases.DepartmentUsecase
}

func NewDepartmentHandler(usecase usecases.DepartmentUsecase) *DepartmentHandler {
	return &DepartmentHandler{
		usecase: usecase,
	}
}

func (h *DepartmentHandler) CreateDept(c *gin.Context) {
	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.CreateDeptRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.CreateDept(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *DepartmentHandler) DeleteDept(c *gin.Context) {
	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.DeleteDeptRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.DeleteDept(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *DepartmentHandler) CreateDeptUser(c *gin.Context) {

	// context 생성
	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.CreateDeptUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.CreateDeptUser(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *DepartmentHandler) DeleteDeptUser(c *gin.Context) {
	// context 생성
	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.DeleteDeptUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.DeleteDeptUser(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}
