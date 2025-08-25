package server

import (
	"encoding/json"
	"org/internal/domains/org/dto/client/requestDTO"
	usecases "org/internal/domains/org/usecases/server"
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

func (h *OrgHandler) CreateOrgFile(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = requestDTO.CreateOrgFileRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.ServerCreateOrgFile(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		// http status code 500
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

	// response.

}
