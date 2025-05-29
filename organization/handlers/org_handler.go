package handlers

import (
	"net/http"
	"org/config"
	"org/usecases"
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

	// response dto 생성

	// request 데이터 파싱 header, body -> dto

	// usecase 호출

	// response.

}

// 눌려서 일부 부서 조회
func (h *OrgHandler) GetDept(w http.ResponseWriter, r *http.Request) {

	// response dto 생성

	// request 데이터 파싱 header, body -> dto

	// usecase 호출

	// response.

}
