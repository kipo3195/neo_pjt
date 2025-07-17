package handlers

import (
	"admin/consts"
	dto "admin/dto/common"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	Org    *OrgHandler
	Common *CommonHandler
}

// 관리자 기본기능 handler

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
