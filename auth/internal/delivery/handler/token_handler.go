package handler

import (
	"auth/internal/application/usecase"
	"auth/internal/consts"
	"auth/internal/delivery/dto/token"
	commonConsts "auth/pkg/consts"
	response "auth/pkg/response"

	"auth/internal/delivery/adapter"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type TokenHandler struct {
	usecase usecase.TokenUsecase
}

func NewTokenHandler(uc usecase.TokenUsecase) *TokenHandler {
	return &TokenHandler{usecase: uc}
}

func (h *TokenHandler) GenerateAppToken(c *gin.Context) {

	ctx := c.Request.Context()

	// ---- 이 로직은 middleware로 처리하기 시작
	// request의 header 데이터 -> dto로 변경
	header := token.GenerateAppTokenRequestHeader{
		Token: c.GetHeader(consts.AUTHORIZATION),
	}

	log.Println("common service에서 호출시 던진 토큰 ", header.Token)

	if header.Token == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	// ---- 이 로직은 middleware로 처리하기 종료

	// request body 데이터 -> dto로 변경
	var req token.GenerateAppTokenRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	input := adapter.MakeGenerateAppTokenInput(req.Uuid)

	// 토큰 발급, DB 저장.
	resDto, err := h.usecase.GenerateAppToken(ctx, input)

	log.Println("handler에서 토큰 구조체 반환 resDto : ", resDto)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, resDto.Body)
	}

	log.Println("handler에서 결과 반환 res : ", resDto.Body)

}

func (h *TokenHandler) AppTokenValidation(c *gin.Context) {

	var body token.AppTokenValidationRequest

	ctx := c.Request.Context()

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	appTokenValidationInput := adapter.MakeAppTokenValidationInput(body.AppToken, body.Token, body.TokenType, body.Uuid)
	result, err := h.usecase.AppTokenValidation(appTokenValidationInput, ctx)

	// 이거 나중에 모듈화 꼭 할 것
	if err != nil || !result { // 에러
		log.Println(err)
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F003, consts.AUTH_F003_MSG)
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *TokenHandler) AppTokenRefresh(c *gin.Context) {

}

func (h *TokenHandler) AccessTokenReIssue(c *gin.Context) {

	ctx := c.Request.Context()

	var req token.AccessTokenReIssueRequest

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

	// appToken 검증
	appTokenValidationInput := adapter.MakeAppTokenValidationInput(req.AppToken, "", "appToken", req.Uuid)
	result, err := h.usecase.AppTokenValidation(appTokenValidationInput, ctx)

	if err != nil {
		// appToken 검증 실패
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F002, consts.AUTH_F002_MSG)
		return
	}

	if result {
		checkRefreshTokenInput := adapter.MakeCheckRefreshTokenInput("", req.Uuid, req.RefreshToken, true)
		result, userId, err := h.usecase.CheckRefreshTokenWithExp(checkRefreshTokenInput, ctx)

		if err != nil {
			if err == consts.ErrUserIdDoesNotExist {
				// uuid : refresh token에 매핑된 사용자 id가 없을때
				response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F011, consts.AUTH_F011_MSG)
			} else {
				// 시간 파싱 에러
				response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
			}
			return
		}

		if result {
			// AT 갱신 처리
			reIssueAccessTokenInput := adapter.MakeReIssueAccessTokenInput(userId, req.Uuid)
			at, err := h.usecase.ReIssueAccessToken(reIssueAccessTokenInput, ctx)
			if err != nil {
				// at 생성 에러
				response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
				return
			}

			// 메모리 로딩을 위한 DB 업데이트 처리
			reIssueAccessTokenSavedInput := adapter.MakeReIssueAccessTokenSavedInput(userId, req.Uuid, req.RefreshToken, at)
			err = h.usecase.ReIssueAccessTokenSaved(ctx, reIssueAccessTokenSavedInput)

			if err != nil {
				// at 생성 에러
				response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
				return
			}

			log.Println("userId : ", userId, " reissued accessToken : ", at)

			res := token.AccessTokenReIssueResponse{
				AccessToken: at,
			}

			response.SendSuccess(c, res)

		} else {
			// 시간 지남 처리
			// device 재등록 처리API 호출 하도록
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.AUTH_F009, consts.AUTH_F009_MSG)
		}

	} else {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

}
