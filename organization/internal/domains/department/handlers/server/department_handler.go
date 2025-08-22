package server

import (
	"encoding/json"
	"net/http"
	"org/consts"
	dto "org/dto/common"
	usecases "org/internal/domains/department/usecases/server"

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
	// response dto 생성
	var res = adminDto.CreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = adminDto.CreateDeptRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data, err := h.usecase.ServerCreateDept(ctx, req)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *DepartmentHandler) DeleteDept(c *gin.Context) {
	// context 생성
	ctx := c.Request.Context()

	// response dto 생성
	var res = adminDto.DeleteDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = adminDto.DeleteDeptRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data, err := h.usecase.ServerDeleteDept(ctx, req)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *DepartmentHandler) CreateDeptUser(c *gin.Context) {

	// context 생성
	// context 생성
	ctx := c.Request.Context()
	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = adminDto.CreateDeptUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data, err := h.usecase.ServerCreateDeptUser(ctx, req)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 500
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusInternalServerError)
	}

	// response.
	json.NewEncoder(w).Encode(res)
}

func (h *DepartmentHandler) DeleteDeptUser(c *gin.Context) {
	// context 생성
	// context 생성
	ctx := c.Request.Context()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = adminDto.DeleteDeptUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data, err := h.usecase.ServerDeleteDeptUser(ctx, req)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 500
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusInternalServerError)
	}

	// response.
	json.NewEncoder(w).Encode(res)
}
