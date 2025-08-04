package server

import (
	"auth/internal/consts"
	"auth/internal/domains/token/dto/server/requestDTO"
	usecases "auth/internal/domains/token/usecases/server"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"

	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	usecase usecases.TokenUsecase
}

func NewTokenHandler(uc usecases.TokenUsecase) *TokenHandler {
	return &TokenHandler{usecase: uc}
}

func (h *TokenHandler) GenerateAppToken(c *gin.Context) {

	// request의 header 데이터 -> dto로 변경
	header := requestDTO.GenerateAppTokenRequestHeader{
		Token: c.GetHeader(consts.AUTHORIZATION),
	}

	log.Println("common service에서 호출시 던진 토큰 ", header.Token)

	if header.Token == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body requestDTO.GenerateAppTokenRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	requestDTO := requestDTO.GenerateAppTokenRequestDTO{
		// Header: header,
		Body: body,
	}

	// 토큰 발급, DB 저장.
	resDto, err := h.usecase.GenerateAppToken(requestDTO.Body)

	log.Println("handler에서 토큰 구조체 반환 resDto : ", resDto)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, resDto.Body)
	}

	log.Println("handler에서 결과 반환 res : ", resDto.Body)

}

func (h *TokenHandler) AppTokenValidation(c *gin.Context) {

	var body requestDTO.AppTokenValidationRequestBody

	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	requestDTO := requestDTO.AppTokenValidationRequestDTO{
		Body: body,
	}

	resDto, err := h.usecase.AppTokenValidation(requestDTO, ctx)

	// 이거 나중에 모듈화 꼭 할 것
	if err != nil || !resDto { // 에러
		log.Println(err)
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *TokenHandler) AppTokenRefresh(c *gin.Context) {

}
