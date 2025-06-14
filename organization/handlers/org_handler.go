package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"org/config"
	"org/consts"
	clDto "org/dto/client"
	dto "org/dto/common"
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

// 눌려서 일부 부서 조회 -> hash가 바뀌었는지 조회 필요.
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
	var res = svDto.SvCreateDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.SvCreateDeptRequest{}

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
	data, err := h.usecase.ServerCreateDept(ctx, req)

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
	var res = svDto.SvDeleteDeptResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.SvDeleteDeptRequest{}

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
	data, err := h.usecase.ServerDeleteDept(ctx, req)

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

func (h *OrgHandler) ServerCreateOrgFile(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.SvCreateOrgFileRequest{}

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

	// usecase 호출
	data, err := h.usecase.ServerCreateOrgFile(ctx, req)

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

func (h *OrgHandler) GetOrgFile(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = clDto.GetOrgFileResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = clDto.GetOrgFileRequest{
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
	data, err := h.usecase.GetOrgFile(ctx, req)

	// response.
	if err == nil {
		// http status code 200
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="org_entity.zip"`)
		w.Write(data.([]byte))
	} else {
		// http status code 400
		res.Code = consts.ERROR
		res.Data = err
		w.WriteHeader(http.StatusBadRequest)
	}

	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *OrgHandler) ServerCreateDeptUser(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.SvCreateDeptUserRequest{}

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

func (h *OrgHandler) ServerDeleteDeptUser(w http.ResponseWriter, r *http.Request) {
	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = svDto.SvDeleteDeptUserRequest{}

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
