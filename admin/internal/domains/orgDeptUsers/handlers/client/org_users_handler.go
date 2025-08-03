package client

import (
	clientUsecase "admin/internal/domains/orgDeptUsers/usecases/client"
	"admin/pkg/consts"
	"encoding/json"

	"admin/internal/domains/orgDeptUsers/dto/client/requestDTO"
	response "admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgDeptUsersHandler struct {
	usecase clientUsecase.OrgDeptUsersUsecase
}

func NewOrgDeptUsersHandler(usecase clientUsecase.OrgDeptUsersUsecase) *OrgDeptUsersHandler {
	return &OrgDeptUsersHandler{
		usecase: usecase,
	}

}

func (h *OrgDeptUsersHandler) RegisterUsers(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body requestDTO.CreateDeptUserRequestBody

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

	requestDTO := requestDTO.CreateDeptUserRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.CreateDeptUser(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *OrgDeptUsersHandler) DeleteUsers(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body requestDTO.DeleteDeptUserRequestBody

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

	requestDTO := requestDTO.DeleteDeptUserRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.DeleteDeptUser(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		response.SendError(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *OrgDeptUsersHandler) GetUsers(c *gin.Context) {

}

func (h *OrgDeptUsersHandler) UpdateUsers(c *gin.Context) {

}
