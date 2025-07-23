package handlers

import (
	consts "auth/consts"
	svCommonReqDto "auth/dto/server/common/request"
	"auth/usecases"
	"auth/utils"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	usecase usecases.ServerUsecase
}

func NewServerHandler(uc usecases.ServerUsecase) *ServerHandler {
	return &ServerHandler{usecase: uc}
}

func (h *ServerHandler) GenerateAppToken(c *gin.Context) {

	// request의 header 데이터 -> dto로 변경
	header := svCommonReqDto.GenerateAppTokenRequestHeader{
		Token: c.GetHeader(consts.AUTHORIZATION),
	}

	log.Println("common service에서 호출시 던진 토큰 ", header.Token)

	if header.Token == "" {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body svCommonReqDto.GenerateAppTokenRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	requestDTO := svCommonReqDto.GenerateAppTokenRequestDTO{
		// Header: header,
		Body: body,
	}

	// 토큰 발급, DB 저장.
	resDto, err := h.usecase.GenerateAppToken(requestDTO.Body)

	log.Println("handler에서 토큰 구조체 반환 resDto : ", resDto)

	if err != nil {
		utils.SendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		utils.SendSuccessResponse(c, resDto.Body)
	}

	log.Println("handler에서 결과 반환 res : ", resDto.Body)

}

func (h *ServerHandler) AppTokenValidation(c *gin.Context) {

	var body svCommonReqDto.AppTokenValidationRequestBody

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	requestDTO := svCommonReqDto.AppTokenValidationRequestDTO{
		Body: body,
	}

	resDto, err := h.usecase.AppTokenValidation(requestDTO, ctx)

	// 이거 나중에 모듈화 꼭 할 것
	if err != nil || !resDto { // 에러
		log.Println(err)
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
	} else {
		utils.SendSuccessResponse(c, "")
	}

}

func (h *ServerHandler) AppTokenRefresh(c *gin.Context) {

}
