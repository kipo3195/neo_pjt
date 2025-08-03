package utils

import (
	consts "admin/pkg/consts"
	"admin/pkg/dto"
	errors "admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, status int, result string, code string, msg string) {

	res := dto.ResponseDTO[errors.ErrorDataDTO]{ // 제네릭 타입 명시 - ResponseDTO의 DATA 'T'에 들어갈 타입을 말함.
		Result: result, // error, fail
		Data: errors.ErrorDataDTO{
			Code:    code,
			Message: msg,
		},
	}
	c.AbortWithStatusJSON(status, res)
}

func SendSuccess[T any](c *gin.Context, t T) {
	res := dto.ResponseDTO[T]{ // 제네릭 타입 명시 - success는 어떤 DTO라도 들어갈 수 있으므로 any
		Result: consts.SUCCESS,
		Data:   t,
	}
	c.AbortWithStatusJSON(200, res) // 200 고정
}
