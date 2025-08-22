package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"org/consts"
	dto "org/dto/common"
	usecases "org/internal/domains/org/usecases/client"

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

// 조직도 전체 조회
func (h *OrgHandler) GetOrgHash(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// response dto 생성
	var res = dto.Response{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.GetOrgHashRequest{
		// 배열의 형태로 받음. org가 하나 이상일 수도 있기 때문.
		OrgHash: r.URL.Query()["orgHash"],
	}

	if len(req.OrgHash) == 0 {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_108,
			Message: consts.E_108_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data, err := h.usecase.GetOrgHash(ctx, req)

	// response.
	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 400
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

func (h *OrgHandler) GetOrgData(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// response dto 생성
	var res = orgDto.GetOrgDataResponse{}

	// request 데이터 파싱 header, body -> dto
	var req = orgDto.GetOrgDataRequest{
		OrgCode: r.URL.Query().Get("orgCode"),
		Type:    r.URL.Query().Get("type"),
		OrgHash: r.URL.Query().Get("orgHash"),
	}

	if len(req.OrgCode) == 0 {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_108,
			Message: consts.E_108_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	// usecase 호출
	file, data, err := h.usecase.GetOrgData(ctx, req)

	// response.
	if file != "" {
		orgCode := req.OrgCode
		// http status code 200
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, orgCode+"_"+file)) // 요청한 org code + 최신 hash
		w.Write(data.([]byte))
		// 전송 헤더의 순서가 영향을 미침 - 파일명 적용이 안됨.
		w.WriteHeader(http.StatusOK)
		return

	} else if err != nil {
		// http status code 400
		w.WriteHeader(http.StatusBadRequest)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	} else {
		res.Result = consts.SUCCESS
		res.Data = data
	}
	// response.
	json.NewEncoder(w).Encode(res)

}
