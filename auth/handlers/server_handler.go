package handlers

import (
	consts "auth/consts"
	dto "auth/dto/common"
	commonSvReqDto "auth/dto/server/common/request"
	"auth/usecases"
	"context"
	"encoding/json"
	"fmt"
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
	header := commonSvReqDto.GenerateAppTokenRequestHeader{
		Token: c.GetHeader(consts.AUTHORIZATION),
	}

	fmt.Println("common service에서 호출시 던진 토큰 ", header.Token)

	if header.Token == "" {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body commonSvReqDto.GenerateAppTokenRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	requestDTO := commonSvReqDto.GenerateAppTokenRequestDTO{
		// Header: header,
		Body: body,
	}

	// 토큰 발급, DB 저장.
	resDto, err := h.usecase.GenerateAppToken(requestDTO.Body)

	fmt.Println("handler에서 토큰 구조체 반환 resDto : ", resDto)

	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	} else {
		sendSuccessResponse(c, resDto.Body)
	}

	fmt.Println("handler에서 결과 반환 res : ", resDto.Body)

}

func sendErrorResponse(c *gin.Context, status int, result string, code string, msg string) {

	res := dto.ResponseDTO[dto.ErrorDataDTO]{ // 제네릭 타입 명시 - ResponseDTO의 DATA 'T'에 들어갈 타입을 말함.
		Result: result, // error, fail
		Data: dto.ErrorDataDTO{
			Code:    code,
			Message: msg,
		},
	}
	c.AbortWithStatusJSON(status, res)
}

func sendSuccessResponse[T any](c *gin.Context, t T) {
	res := dto.ResponseDTO[T]{ // 제네릭 타입 명시 - success는 어떤 DTO라도 들어갈 수 있으므로 any
		Result: consts.SUCCESS,
		Data:   t,
	}
	c.AbortWithStatusJSON(200, res) // 200 고정
}

func (h *ServerHandler) AppTokenValidation(c *gin.Context) {

	var body commonSvReqDto.AppTokenValidationRequestBody

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	requestDTO := commonSvReqDto.AppTokenValidationRequestDTO{
		Body: body,
	}

	resDto, err := h.usecase.AppTokenValidation(requestDTO, ctx)

	// 이거 나중에 모듈화 꼭 할 것
	if err != nil || !resDto { // 에러
		fmt.Println(err)
		sendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
	} else {
		sendSuccessResponse(c, "")
	}

}
