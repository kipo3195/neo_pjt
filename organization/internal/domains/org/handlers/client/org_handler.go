package client

import (
	"org/internal/domains/org/dto/client/requestDTO"
	usecases "org/internal/domains/org/usecases/client"
	commonConsts "org/pkg/consts"
	"org/pkg/response"

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

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.GetOrgHashRequest{
		// 배열의 형태로 받음. org가 하나 이상일 수도 있기 때문.
		OrgHash: c.QueryArray("orgHash"),
	}

	if len(req.OrgHash) == 0 {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.GetOrgHash(ctx, req)

	// response.
	if err == nil {
		response.SendSuccess(c, data)
	} else {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}
}

func (h *OrgHandler) GetOrgData(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.GetOrgDataRequest{
		OrgCode: c.Query("orgCode"),
		Type:    c.Query("type"),
		OrgHash: c.Query("orgHash"),
	}

	if len(req.OrgCode) == 0 {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}
	// usecase 호출
	file, data, err := h.usecase.GetOrgData(ctx, req)

	// response.
	if file != "" {
		orgCode := req.OrgCode
		// http status code 200
		// w.Header().Set("Content-Type", "application/octet-stream")
		// w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, orgCode+"_"+file)) // 요청한 org code + 최신 hash
		// w.Write(data.([]byte))
		// // 전송 헤더의 순서가 영향을 미침 - 파일명 적용이 안됨.
		// w.WriteHeader(http.StatusOK)
		response.SendFileStream(c, data, orgCode+"_"+file+".zip", "")
		return

	} else if err != nil {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, data)
	}
}
