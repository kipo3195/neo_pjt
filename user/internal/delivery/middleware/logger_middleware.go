package middleware

import (
	"context"
	"time"
	"user/internal/domain/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()

		// 1. TraceID를 표준 context에 주입
		ctx := context.WithValue(c.Request.Context(), "trace_id", traceID)

		// 2. Gin Context의 Request를 새로운 context로 교체
		c.Request = c.Request.WithContext(ctx)

		// defer를 통해 모든 핸들러가 끝난 후 로그 출력
		defer func() {

			// c.Next()를 거치며 미들웨어가 추가하더라도 그 미들웨어에서 업데이트한 데이터를 꺼내 오기 위함
			latestCtx := c.Request.Context()

			// 성공이든 실패든 무조건 실행되는 '최종 결과 요약'
			// 상태 코드가 400, 500이면 실패인 걸 이미 status 필드가 말해주고 있습니다.
			logger.Info(latestCtx, "request_finished",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(), // 여기서 200인지 500인지 찍힘
				"latency", time.Since(start),
			)
		}()

		// 다음 미들웨어 또는 핸들러 실행
		c.Next()
	}
}
