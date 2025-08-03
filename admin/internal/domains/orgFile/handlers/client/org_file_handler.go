package client

import (
	requestDTO "admin/internal/domains/orgFile/dto/client/requestDTO"
	clientUsecase "admin/internal/domains/orgFile/usecases/client"
	response "admin/pkg/response"
	"encoding/json"

	commonConsts "admin/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgFileHandler struct {
	usecase clientUsecase.OrgFileUsecase
}

func NewOrgFileHandler(usecase clientUsecase.OrgFileUsecase) *OrgFileHandler {

	return &OrgFileHandler{
		usecase: usecase,
	}

}

func (h *OrgFileHandler) CreateOrgFile(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body requestDTO.CreateOrgFileRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	requestDTO := requestDTO.CreateOrgFileRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.CreateOrgFile(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}
}

func (h *OrgFileHandler) GetOrgFile(c *gin.Context) {

}
