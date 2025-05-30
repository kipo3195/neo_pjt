package handlers

import (
	"admin/consts"
	"admin/dto"
	"admin/usecases"
	"encoding/json"
	"net/http"
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
	var res = dto.CreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = dto.CreateDeptRequest{}

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

	// usecase 호출
	data, err := h.usecase.CreateDepartment(ctx, req)

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

}
