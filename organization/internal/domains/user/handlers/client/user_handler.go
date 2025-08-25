package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"org/consts"
	"org/contextkey"
	dto "org/dto/common"
	usecases "org/internal/domains/user/usecases/client"
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

	// response dto 생성
	var res = dto.Response{}

	// 인증 토큰에서 요청 사용자의 hash 정보 추출
	val := r.Context().Value(contextkey.UserHashKey)
	myHash, ok := val.(string)
	if !ok {
		log.Println("인증 토큰의 userHash 데이터 에러 ")
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.ORG_F101,
			Message: consts.ORG_F101_MSG,
		}
		return
	}

	log.Println("내 정보 요청시 myHash : ", myHash)
	// dto 생성
	var req = userDto.GetMyInfoRequest{
		MyHash: myHash,
	}

	data, err := h.usecase.GetMyInfo(ctx, req)

	// response.
	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// http status code 400
		w.WriteHeader(http.StatusInternalServerError)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}
	// response.
	json.NewEncoder(w).Encode(res)

}

func (h *UserHandler) GetUserInfo(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

}
