package handler

import (
	"encoding/json"
	"log"
	"org/internal/application/orchestrator"
	"org/internal/delivery/adapter"
	"org/internal/delivery/dto/dummy"
	commonConsts "org/pkg/consts"
	"org/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type DummyDataServiceHandler struct {
	svc *orchestrator.DummyDataService
}

func NewDummyDataServiceHandler(svc *orchestrator.DummyDataService) *DummyDataServiceHandler {
	return &DummyDataServiceHandler{svc: svc}
}

func (h *DummyDataServiceHandler) InitServiceUser(c *gin.Context) {
	ctx := c.Request.Context()
	//.. init service 비즈니스 로직 작성

	var req dummy.CreateServiceUserRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeCreateServiceUserInput(req.UserCount, req.Keyword)
	err := h.svc.User.CreateServiceUser(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}

// func (h *DummyDataServiceHandler) InitUserDetail(c *gin.Context) {

// 	ctx := c.Request.Context()

// 	var req dummy.CreateUserDetailRequest

// 	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
// 		return
// 	}

// 	// 필수 데이터 검증
// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
// 		return
// 	}

// 	input := adapter.MakeCreateUserDetailInput(req.Keyword, req.Type)
// 	err := h.svc.User.CreateUserDetail(ctx, input)

// 	if err != nil {
// 		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
// 	} else {
// 		response.SendSuccess(c, "success")
// 	}

// }

func (h *DummyDataServiceHandler) InitUserMultiLang(c *gin.Context) {

	ctx := c.Request.Context()

	var req dummy.CreateUserMultiLangRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeUserMultiLangInput(req.Keyword)
	err := h.svc.User.CreateUserMultiLang(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}

func (h *DummyDataServiceHandler) InitWorksDept(c *gin.Context) {
	ctx := c.Request.Context()

	var req dummy.CreateWorksDeptRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeWorksDeptInput(req.Org, req.MaxDepth, req.DeptCount)
	err := h.svc.Department.CreateWorksDept(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}

func (h *DummyDataServiceHandler) InitWorksDeptMultiLang(c *gin.Context) {

	ctx := c.Request.Context()

	var req dummy.CreateWorksDeptMultiLangRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeCreateWorksDeptMultiLangInput(req.Org)
	err := h.svc.Department.CreateWorksDeptMultiLang(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}

func (h *DummyDataServiceHandler) InitWorksDeptUser(c *gin.Context) {

	ctx := c.Request.Context()

	var req dummy.CreateWorksDeptUserRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeCreateWorksDeptUserInput(req.Org)

	output, err := h.svc.User.GetServiceUsers(ctx, input.Org)

	if err != nil {
		log.Println("getService user 조회 에러 11")
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	if output == nil || len(output) == 0 {
		log.Println("getService user 조회 에러 22")
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	temp := make([]string, len(output))
	for i := 0; i < len(output); i++ {
		temp[i] = output[i].UserHash
	}

	input.UserHashs = temp
	err = h.svc.Department.CreateWorksDeptUser(ctx, input)

}
