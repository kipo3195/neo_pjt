package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		// context.WithTimeout으로 타임아웃 설정 (context.WithTimeout에서 발생하는 "신호"는 명시적 호출 없이 자동으로 발생합니다.)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
	}
}
