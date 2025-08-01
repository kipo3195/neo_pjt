package handlers

import (
	"auth/internal/consts"
	"auth/internal/domains/certification/dto/client/requestDTO"
	clientUsecases "auth/internal/domains/certification/usecases/client"
	usecases "auth/internal/domains/certification/usecases/client"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"
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
	header := requestDTO.AuthRequestHeader{
		Token: c.GetHeader(consts.LOGIN_HEADER_AUTH_TOKEN),
		Uuid:  c.GetHeader(consts.LOGIN_HEADER_UUID),
	}

	// header 검증
	if header.Token == "" || header.Uuid == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	// request body 데이터 -> dto로 변경
	var body requestDTO.AuthRequestBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	requestDTO := requestDTO.AuthRequestDTO{
		Header: header,
		Body:   body,
	}

	// 비즈니스 로직 호출
	resDto, err := h.usecase.GetAuth(requestDTO)

	if err != nil {

		if err == consts.ErrUnregisteredUuid {
			// 등록된 UUID가 아님
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F001, consts.AUTH_F001_MSG)
		} else if err == consts.ErrTokenMismatch {
			// 토큰 정보 불일치, 재발급 필요.
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F002, consts.AUTH_F002_MSG)
		} else if err == consts.ErrAuthenticationFailed {
			// 등록된 사용자가 없음.
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
		} else if err == consts.ErrUnregisteredUser {
			// 등록된 사용자가 아님.
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F004, consts.AUTH_F004_MSG)
		} else {
			// server error : db, jwt make
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}

	} else {
		response.SendSuccess(c, resDto.Body)
	}
}
