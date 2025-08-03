package client

import (
	clientUsecase "admin/internal/domains/orgDepts/usecases/client"
	"admin/pkg/consts"
	"encoding/json"

	"admin/internal/domains/orgDepts/dto/client/requestDTO"
	response "admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgDeptsHandler struct {
	usecase clientUsecase.OrgDeptsUsecase
}

func NewOrgDeptsHandler(usecase clientUsecase.OrgDeptsUsecase) *OrgDeptsHandler {
	return &OrgDeptsHandler{
		usecase: usecase,
	}

}

func (h *OrgDeptsHandler) RegisterDepts(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	// 이미 timeout이 걸려있는 ctx이므로 그대로 사용만 하면됨.
	ctx := c.Request.Context()

	var body requestDTO.RegisterDeptRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := requestDTO.RegisterDeptRequestDTO{
		Body: body,
	}

	// 필요시 result 값 response status로 분기 처리
	_, err := h.usecase.CreateDepartment(ctx, requestDTO)

	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}
}

func (h *OrgDeptsHandler) GetDepts(c *gin.Context) {

}

func (h *OrgDeptsHandler) UpdateDepts(c *gin.Context) {

}

func (h *OrgDeptsHandler) DeleteDepts(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body requestDTO.DeleteDeptRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		response.SendError(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := requestDTO.DeleteDeptRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.DeleteDepartment(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}
