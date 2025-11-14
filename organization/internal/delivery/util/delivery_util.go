package util

import (
	"org/internal/consts"

	"github.com/gin-gonic/gin"
)

func GetUserIdByAccessToken(c *gin.Context) string {

	temp := c.Value(consts.USER_ID)
	userId, ok := temp.(string)
	if !ok {
		return ""
	} else {
		return userId
	}
}
func GetUserHashByAccessToken(c *gin.Context) string {

	temp := c.Value(consts.USER_HASH)
	userHash, ok := temp.(string)
	if !ok {
		return ""
	} else {
		return userHash
	}
}
