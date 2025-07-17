package handlers

import (
	"admin/consts"
	orgDto "admin/dto/client/org"
	clOrgReqDto "admin/dto/client/org/request"
	dto "admin/dto/common"
	"admin/usecases"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 관리자 - org 서비스 연계 handler
type OrgHandler struct {
	usecase usecases.OrgUsecase
}

func NewOrgHandler(r usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: r}
}

func (h *OrgHandler) CreateDept(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	// 이미 timeout이 걸려있는 ctx이므로 그대로 사용만 하면됨.
	ctx := c.Request.Context()

	var body clOrgReqDto.CreateDeptRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.CreateDeptRequestDTO{
		Body: body,
	}

	// usecase 호출
	err := h.usecase.CreateDepartment(ctx, requestDTO)

	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}
}

func (h *OrgHandler) GetDept(w http.ResponseWriter, r *http.Request) {

}

func (h *OrgHandler) UpdateDept(w http.ResponseWriter, r *http.Request) {

}

func (h *OrgHandler) DeleteDept(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = orgDto.DeleteDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.DeleteDeptRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// usecase 호출
	data, err := h.usecase.DeleteDepartment(ctx, req)

	fmt.Println(data)

	if err == nil {
		// http status code 200
		res.Code = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 400
		res.Code = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *OrgHandler) CreateOrgFile(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = orgDto.CreateOrgFileResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.CreateOrgFileRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// usecase 호출
	data, err := h.usecase.CreateOrgFile(ctx, req)

	fmt.Println(data)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// 서버 - 서버 통신이 실패했다는 의미.
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusInternalServerError)
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *OrgHandler) GetOrgFile(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	//ctx := r.Context()

	// response dto 생성
	var res = orgDto.CreateOrgFileResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.CreateOrgFileRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
}

func (h *OrgHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = orgDto.CreateDeptUserResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.CreateDeptUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// usecase 호출
	data, err := h.usecase.CreateDeptUser(ctx, req)

	fmt.Println(data)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 400
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *OrgHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.DeleteDeptUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// usecase 호출
	data, err := h.usecase.DeleteDeptUser(ctx, req)

	fmt.Println(data)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 400
		res.Result = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *OrgHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (h *OrgHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}
