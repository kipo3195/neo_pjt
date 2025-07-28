package handlers

import (
	"auth/internal/consts"
	clReqDto "auth/internal/domains/certification/dto/client/request"
	clientUsecases "auth/internal/domains/certification/usecases/client"
	usecases "auth/internal/domains/certification/usecases/client"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type CertificationHandler struct {
	usecase clientUsecases.CertificationUsecase
}

func NewCertificationHandler(uc usecases.CertificationUsecase) *CertificationHandler {
	return &CertificationHandler{usecase: uc}
}

func (h *CertificationHandler) Login(c *gin.Context) {

	// request의 header 데이터 -> dto로 변경
	header := clReqDto.AuthRequestHeader{
		Token: c.GetHeader(consts.LOGIN_HEADER_AUTH_TOKEN),
		Uuid:  c.GetHeader(consts.LOGIN_HEADER_UUID),
	}

	// header 검증
	if header.Token == "" || header.Uuid == "" {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body clReqDto.AuthRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_103, consts.E_103_MSG)
		return
	}

	requestDTO := clReqDto.AuthRequestDTO{
		Header: header,
		Body:   body,
	}

	// 비즈니스 로직 호출
	resDto, err := h.usecase.GetAuth(requestDTO)

	if err != nil {

		if err == consts.ErrUnregisteredUuid {
			// 등록된 UUID가 아님
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F001, consts.AUTH_F001_MSG)
		} else if err == consts.ErrTokenMismatch {
			// 토큰 정보 불일치, 재발급 필요.
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F002, consts.AUTH_F002_MSG)
		} else if err == consts.ErrAuthenticationFailed {
			// 등록된 사용자가 없음.
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
		} else if err == consts.ErrUnregisteredUser {
			// 등록된 사용자가 아님.
			utils.SendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.AUTH_F004, consts.AUTH_F004_MSG)
		} else {
			// server error : db, jwt make
			utils.SendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}

	} else {
		utils.SendSuccessResponse(c, resDto.Body)
	}
}
