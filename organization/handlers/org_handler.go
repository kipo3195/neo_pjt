package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"org/config"
	"org/consts"
	"org/dto"
	clDto "org/dto/client"
	svDto "org/dto/server"
	"org/usecases"
	"time"
)

type OrgHandler struct {
	usecase usecases.OrgUsecase
	sfg     *config.ServerConfig
}

func NewOrgHandler(sfg *config.ServerConfig, uc usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: uc, sfg: sfg}
}

// 조직도 전체 조회
func (h *OrgHandler) GetOrg(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = clDto.GetOrgResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.GetOrgRequest{
		OrgCode: r.URL.Query().Get("orgCode"),
	}

	if req.OrgCode == "" {
		res.Code = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_104,
			Message: consts.E_104_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	// usecase 호출
	data, err := h.usecase.GetOrg(ctx, req)

	// response.
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

// 눌려서 일부 부서 조회
func (h *OrgHandler) GetDept(w http.ResponseWriter, r *http.Request) {

	// context 생성

	// response dto 생성

	// request 데이터 파싱 header, body -> dto

	// usecase 호출

	// response.

}

func (h *OrgHandler) ServerCreateDept(w http.ResponseWriter, r *http.Request) {
	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = svDto.ServerCreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.ServerCreateDeptRequest{}

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
	data, err := h.usecase.ServerCreateDepartment(ctx, req)

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

func (h *OrgHandler) ServerDeleteDept(w http.ResponseWriter, r *http.Request) {
	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = svDto.ServerDeleteDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.ServerDeleteDeptRequest{}

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
	data, err := h.usecase.ServerDeleteDepartment(ctx, req)

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
