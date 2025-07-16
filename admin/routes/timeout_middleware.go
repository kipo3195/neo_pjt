package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		// context.WithTimeout으로 타임아웃 설정 (context.WithTimeout에서 발생하는 "신호"는 명시적 호출 없이 자동으로 발생합니다.)
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		// defer 없이 cancel() // 수동으로 취소

		// 새 context로 request 갱신
		c.Request = c.Request.WithContext(ctx)

		// 채널을 이용해 완료 여부 체크
		finished := make(chan struct{})
		go func() {
			c.Next() // 다음 미들웨어/핸들러 실행
			close(finished)
		}()

		select {
		case <-ctx.Done():
			// 타임아웃 발생
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timed out",
			})
		case <-finished:
			// 정상 완료
		}
	}
}
