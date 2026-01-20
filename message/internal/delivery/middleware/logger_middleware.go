package middleware

import (
	"context"
	"message/internal/domain/logger"
	"time"

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
			logger.Info(ctx, "request_completed",
				"trace_id", traceID,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(), // 응답 상태 코드
				"latency", time.Since(start),
			)
		}()

		// 다음 미들웨어 또는 핸들러 실행
		c.Next()
	}
}
