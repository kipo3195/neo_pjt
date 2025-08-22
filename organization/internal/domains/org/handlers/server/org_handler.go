package server

import (
	"encoding/json"
	"net/http"
	"org/consts"
	dto "org/dto/common"
	usecases "org/internal/domains/org/usecases/server"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	usecase usecases.OrgUsecase
}

func NewOrgHandler(usecase usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{
		usecase: usecase,
	}
}

func (h *OrgHandler) CreateOrgFile(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = adminDto.CreateOrgFileRequest{}

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
