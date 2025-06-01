package handlers

import (
	"admin/consts"
	clDto "admin/dto/client"
	dto "admin/dto/common"
	"admin/usecases"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AdminHandler struct {
	usecase usecases.AdminUsecase
}

func NewAdminHandler(r usecases.AdminUsecase) *AdminHandler {
	return &AdminHandler{usecase: r}
}

func (h *AdminHandler) CreateDept(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()

	// response dto 생성
	var res = clDto.CreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.CreateDeptRequest{}

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
	data, err := h.usecase.CreateDepartment(ctx, req)

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

func (h *AdminHandler) GetDept(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminHandler) UpdateDept(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminHandler) DeleteDept(w http.ResponseWriter, r *http.Request) {

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
