package handlers

import (
	"admin/consts"
	clOrgReqDto "admin/dto/client/org/request"
	"admin/usecases"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 관리자 - org 서비스 연계 handler
type OrgHandler struct {
	usecase usecases.OrgUsecase
}

func NewOrgHandler(r usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: r}
}

func (h *OrgHandler) CreateDept(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	// 이미 timeout이 걸려있는 ctx이므로 그대로 사용만 하면됨.
	ctx := c.Request.Context()

	var body clOrgReqDto.CreateDeptRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.CreateDeptRequestDTO{
		Body: body,
	}

	// 필요시 result 값 response status로 분기 처리
	_, err := h.usecase.CreateDepartment(ctx, requestDTO)

	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}
}

func (h *OrgHandler) GetDept(c *gin.Context) {

}

func (h *OrgHandler) UpdateDept(c *gin.Context) {

}

func (h *OrgHandler) DeleteDept(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body clOrgReqDto.DeleteDeptRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.DeleteDeptRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.DeleteDepartment(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}

}

func (h *OrgHandler) CreateOrgFile(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body clOrgReqDto.CreateOrgFileRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.CreateOrgFileRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.CreateOrgFile(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}
}

func (h *OrgHandler) GetOrgFile(c *gin.Context) {

}

func (h *OrgHandler) CreateUser(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body clOrgReqDto.CreateDeptUserRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.CreateDeptUserRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.CreateDeptUser(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}

}

func (h *OrgHandler) DeleteUser(c *gin.Context) {
	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var body clOrgReqDto.DeleteDeptUserRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_108, consts.E_108_MSG)
		return
	}

	requestDTO := clOrgReqDto.DeleteDeptUserRequestDTO{
		Body: body,
	}

	// usecase 호출
	_, err := h.usecase.DeleteDeptUser(ctx, requestDTO)

	// 필요시 result 값 response status로 분기 처리
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, "")
	}

}

func (h *OrgHandler) GetUser(c *gin.Context) {

}

func (h *OrgHandler) UpdateUser(c *gin.Context) {

}
