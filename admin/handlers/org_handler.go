package handlers

import (
	"admin/consts"
	clDto "admin/dto/client"
	dto "admin/dto/common"
	"admin/usecases"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type OrgHandler struct {
	usecase usecases.AdminOrgUsecase
}

func NewAdminHandler(r usecases.AdminOrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: r}
}

func (h *OrgHandler) CreateDept(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = clDto.CreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.CreateDeptRequest{}

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
		w.WriteHeader(http.StatusBadRequest)
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_108,
			Message: consts.E_108_MSG,
		}
		return
	}

	// usecase 호출
	data, err := h.usecase.CreateDepartment(ctx, req)

	fmt.Println(data)

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

func (h *OrgHandler) GetDept(w http.ResponseWriter, r *http.Request) {

}

func (h *OrgHandler) UpdateDept(w http.ResponseWriter, r *http.Request) {

}

func (h *OrgHandler) DeleteDept(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = clDto.DeleteDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.DeleteDeptRequest{}

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
	var res = clDto.CreateOrgFileResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.CreateOrgFileRequest{}

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
	var res = clDto.CreateOrgFileResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.CreateOrgFileRequest{}

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
	var res = clDto.CreateDeptUserResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.CreateDeptUserRequest{}

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
	var req = clDto.DeleteDeptUserRequest{}

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
