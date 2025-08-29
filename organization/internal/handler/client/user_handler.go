package client

import (
	"context"
	"log"
	"org/internal/consts"
	"org/internal/domains/user/dto/client/requestDTO"
	usecases "org/internal/domains/user/usecases/client"
	"org/internal/middleware/contextkey"
	commonConsts "org/pkg/consts"
	"org/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(usecase usecases.UserUsecase) *UserHandler {

	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) GetMyInfo(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// 인증 토큰에서 요청 사용자의 hash 정보 추출
	val := c.Value(contextkey.UserHashKey)
	myHash, ok := val.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, consts.ORG_F101, consts.ORG_F101_MSG)
		log.Println("인증 토큰의 userHash 데이터 에러 ")
		return
	}

	log.Println("내 정보 요청시 myHash : ", myHash)
	// dto 생성
	var req = requestDTO.GetMyInfoRequest{
		MyHash: myHash,
	}

	data, err := h.usecase.GetMyInfo(ctx, req)

	// response.
	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *UserHandler) GetUserInfo(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

}
